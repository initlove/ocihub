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
}

func (m *Memory) Init(paras map[string]interface{}) error {
	m.store = make(map[string]session.Record)
	return nil
}

func (m *Memory) New(ctx context.Context) (string, error) {
	if m.store == nil {
		return "", errors.New("Please init the 'session memory' driver before use it.")
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
	}
}
