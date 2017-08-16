package session

import (
	"errors"
	"fmt"
	"sync"

	"github.com/astaxie/beego/context"
)

type SessionDriver interface {
	Init(paras map[string]interface{}) error
	New(ctx context.Context) (string, error)
	Get(ctx context.Context, id string) (interface{}, error)
	Release(ctx context.Context, id string) error
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

func InitSession(name string, paras map[string]interface{}) error {
	for n, d := range sds {
		if name == n {
			err := d.Init(paras)
			if err == nil {
				sysSession = d
			}
			return err
		}
	}

	return fmt.Errorf("SessionDriver '%s' is not supported.", name)
}

func New(ctx context.Context) (string, error) {
	if sysSession == nil {
		return "", errors.New("Please init the session driver first.")
	}

	return sysSession.New(ctx)
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

func GC() error {
	if sysSession == nil {
		return errors.New("Please init the session driver first.")
	}

	return sysSession.GC()
}
