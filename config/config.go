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
)

type Config struct {
	Port     int64    `yaml:"port"`
	LogLevel string   `yaml:"loglevel,omitempty"`
	DB       DBConfig `yaml:"db"`
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

var (
	sysConfig Config
)

func LoadConfigFile(path string) (Config, error) {
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
