package session

import (
	"errors"
	"fmt"
	"sync"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"

	"github.com/initlove/ocihub/config"
)

type SessionDriver interface {
	Init(paras map[string]interface{}) error
	New(ctx context.Context, id string) (string, error)
	Get(ctx context.Context, id string) (interface{}, error)
	Release(ctx context.Context, id string) error

	// Should merge to 'Get' and add 'name' paras
	GetCache(ctx context.Context, id string) ([]byte, error)
	PutCache(ctx context.Context, id string, data []byte) error
	GC() error
}

var (
	sdLock sync.Mutex
	sds    = make(map[string]SessionDriver, 16)

	sysSession SessionDriver = nil
)

func Register(name string, driver SessionDriver) error {
	if name == "" {
		return errors.New("Could not register a session driver with empty name.")
	}

	if driver == nil {
		return errors.New("Could not register a nil session driver.")
	}

	sdLock.Lock()
	defer sdLock.Unlock()

	if _, exists := sds[name]; exists {
		return fmt.Errorf("SessionDriver '%s' is already registered.", name)
	}

	sds[name] = driver

	return nil
}

func InitSession(cfg config.SessionConfig) error {
	for n, v := range cfg {
		if d, ok := sds[n]; ok {
			logs.Debug("Init Session Driver: '%s'.", n)
			err := d.Init(v)
			if err == nil {
				sysSession = d
			}
			return err
		}
	}

	return errors.New("Cannot find supported session driver.")
}

func New(ctx context.Context, id string) (string, error) {
	if sysSession == nil {
		return "", errors.New("Please init the session driver first.")
	}

	return sysSession.New(ctx, id)
}

func Get(ctx context.Context, id string) (interface{}, error) {
	if sysSession == nil {
		return nil, errors.New("Please init the session driver first.")
	}

	return sysSession.Get(ctx, id)
}

func Release(ctx context.Context, id string) error {
	if sysSession == nil {
		return errors.New("Please init the session driver first.")
	}

	return sysSession.Release(ctx, id)
}

func GetCache(ctx context.Context, id string) ([]byte, error) {
	if sysSession == nil {
		return nil, errors.New("Please init the session driver first.")
	}

	// FIXME: id should not be printed after service online
	logs.Debug("Session GetCache '%s'.", id)
	return sysSession.GetCache(ctx, id)
}

func PutCache(ctx context.Context, id string, data []byte) error {
	if sysSession == nil {
		return errors.New("Please init the session driver first.")
	}

	// FIXME: id should not be printed after service online
	logs.Debug("Session PutCache '%s'.", id)
	return sysSession.PutCache(ctx, id, data)
}

func GC() error {
	if sysSession == nil {
		return errors.New("Please init the session driver first.")
	}

	return sysSession.GC()
}
