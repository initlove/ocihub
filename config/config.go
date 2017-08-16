package config

import (
	"errors"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var (
	EMPTY_DB_USER_OR_PASSWD = errors.New("The database 'User' or 'Password' should not be empty.")
	EMPTY_DB_NAME           = errors.New("The database 'Name' should not be empty.")
	EMPTY_DB_SERVER         = errors.New("The database 'Server' should not be empty.")

	NON_STORAGE_BACKEND = errors.New("At least one storage backend required.")
	NON_SESSION_BACKEND = errors.New("At least one session backend required.")
)

type Config struct {
	Port int64 `yaml:"port"`
	// The log validation will be checked in the log init part.
	Log LogConfig `yaml:"log,omitempty"`
	// 'dynamic' or 'static'
	// static: load at the first time
	// dynamic: load every time, most time because of multiple tenant using their own token/ak-sk
	// TODO: should have 'default' value
	StorageLoad string        `yaml:"storageload,omitempty"`
	Storage     StorageConfig `yaml:"storage"`
	DB          DBConfig      `yaml:"db"`
	Session     SessionConfig `yaml:"session"`
}

func (cfg *Config) Valid() error {
	if err := cfg.Storage.Valid(); err != nil {
		return err
	}
	if err := cfg.DB.Valid(); err != nil {
		return err
	}
	if err := cfg.Session.Valid(); err != nil {
		return err
	}

	return nil
}

type LogConfig map[string](map[string]interface{})

type StorageConfig map[string](map[string]interface{})

func (cfg *StorageConfig) Valid() error {
	if len(*cfg) == 0 {
		return NON_STORAGE_BACKEND
	}

	return nil
}

type SessionConfig map[string](map[string]interface{})

func (cfg *SessionConfig) Valid() error {
	if len(*cfg) == 0 {
		return NON_SESSION_BACKEND
	}

	return nil
}

type DBConfig struct {
	Driver   string `yaml:"driver"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Server   string `yaml:"server"`
	Name     string `yaml:"name"`
}

func (cfg *DBConfig) GetConnection() (string, error) {
	if cfg.User == "" || cfg.Password == "" {
		return "", EMPTY_DB_USER_OR_PASSWD
	}

	if cfg.Name == "" {
		return "", EMPTY_DB_NAME
	}
	if cfg.Server == "" {
		return "", EMPTY_DB_SERVER
	}

	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", cfg.User, cfg.Password, cfg.Server, cfg.Name), nil
}

func (cfg *DBConfig) Valid() error {
	_, err := cfg.GetConnection()
	if err != nil {
		return err
	}

	return nil
}

var (
	sysConfig Config
)

func InitConfigFromFile(path string) (Config, error) {
	var config Config
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	// TODO: add lock?
	sysConfig = config

	return config, nil
}

func GetConfig() Config {
	return sysConfig
}
