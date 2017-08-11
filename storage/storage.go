package storage

import (
	"errors"
	"fmt"

	"github.com/initlove/ocihub/config"
	"github.com/initlove/ocihub/health"
	"github.com/initlove/ocihub/storage/driver"
)

type StorageHealth struct {
}

func (sh *StorageHealth) GetStatus() (string, string) {
	return "", ""
}

func HealthCheck() error {
	health.RegisterHealth("storage", &StorageHealth{})
	return nil
}

var (
	sysDriver driver.StorageDriver = nil
)

// TODO: more logs
func loadDriver() (driver.StorageDriver, error) {
	cfg := config.GetConfig().Storage
	for n, paras := range cfg {
		d, err := driver.FindDriver(n, paras)
		if err == nil {
			// Pickup the first qualified driver
			err = d.Create(paras)
			if err == nil {
				return d, nil
			}
		}
	}

	return nil, errors.New("Fail to get a suitable storage driver")
}

func InitStorage() error {
	_, err := loadDriver()
	// TODO: we should check the healthy status at the beginning
	return err
}

func Driver() driver.StorageDriver {
	if config.GetConfig().StorageLoad == "static" && sysDriver != nil {
		return sysDriver
	}

	var err error
	sysDriver, err = loadDriver()
	if err != nil {
		panic("Failed to load driver")
	}

	return sysDriver
}
