package router

import (
	"errors"
	"fmt"
	"sync"

	"github.com/astaxie/beego"
)

var (
	nssLock sync.Mutex
	nss     = make(map[string]*beego.Namespace, 8)
)

func RegisterRouter(name string, ns *beego.Namespace) error {
	if name == "" {
		return errors.New("Cannot register a namespace without a name")
	}

	if ns == nil {
		return errors.New("Cannot register a nil namespace")
	}

	nssLock.Lock()
	defer nssLock.Unlock()

	if _, existed := nss[name]; existed {
		return fmt.Errorf("Namespace '%s' is already exist.", name)
	}
	nss[name] = ns

	return nil
}

func GetNamespaces() map[string]*beego.Namespace {
	return nss
}
