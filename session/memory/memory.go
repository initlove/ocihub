package memory

import (
	"errors"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/satori/go.uuid"

	"github.com/initlove/ocihub/session"
)

const (
	memoryPrefix = "inmemory"
)

type Memory struct {
	store map[string]session.Record
	cache map[string][]byte
}

func (m *Memory) Init(paras map[string]interface{}) error {
	m.store = make(map[string]session.Record)
	m.cache = make(map[string][]byte)
	return nil
}

func (m *Memory) New(ctx context.Context, id string) (string, error) {
	if m.store == nil {
		return "", errors.New("Please init the 'session memory' driver before use it.")
	}

	if id != "" {
		//FIXME: check id
		m.store[id] = session.NewRecordFromContext(ctx)
		return id, nil
	}

	sessionUUID := uuid.NewV4().String()
	m.store[sessionUUID] = session.NewRecordFromContext(ctx)
	return sessionUUID, nil
}

//TODO: the return interface is not designed, useless?
func (m *Memory) Get(ctx context.Context, id string) (interface{}, error) {
	if m.store == nil {
		return nil, errors.New("Please init the 'session memory' driver before use it.")
	}

	r, ok := m.store[id]
	if !ok {
		return nil, errors.New("Cannot get the matched sessionid")
	}

	err := r.Match(ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (m *Memory) Release(ctx context.Context, id string) error {
	if m.store == nil {
		return errors.New("Please init the 'session memory' driver before use it.")
	}

	if _, ok := m.store[id]; !ok {
		return errors.New("Cannot get the matched sessionid")
	}

	delete(m.store, id)
	delete(m.cache, id)
	return nil
}

//TODO: id should be same with session id
func (m *Memory) GetCache(ctx context.Context, id string) ([]byte, error) {
	if m.cache == nil {
		return nil, errors.New("Please init the 'session memory' driver before use it.")
	}

	data, ok := m.cache[id]
	if !ok {
		return nil, errors.New("Cannot get the matched cache data.")
	}
	return data, nil
}

func (m *Memory) PutCache(ctx context.Context, id string, data []byte) error {
	if m.cache == nil {
		return errors.New("Please init the 'session memory' driver before use it.")
	}

	m.cache[id] = data
	return nil
}

func (m *Memory) GC() error {
	if m.store == nil {
		return errors.New("Please init the 'session memory' driver before use it.")
	}

	var expired []string
	for k, r := range m.store {
		if r.Expired() {
			expired = append(expired, k)
		}
	}

	num := len(expired)
	if num == 0 {
		return nil
	}

	logs.Info("Session GC start: '%d' expired session detected.", num)
	for _, id := range expired {
		delete(m.store, id)
	}

	return nil
}

func init() {
	if err := session.Register(memoryPrefix, &Memory{}); err != nil {
		logs.Error("Failed to register memory session driver.")
	} else {
		logs.Debug("Session driver '%s' registered.", memoryPrefix)
	}
}
