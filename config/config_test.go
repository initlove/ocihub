package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConnection(t *testing.T) {
	cases := []struct {
		cfg  DBConfig
		conn string
		err  error
	}{
		{DBConfig{"mysql", "user", "passwd", "server", "name"}, "user:passwd@tcp(server)/name?charset=utf8", nil},
		{DBConfig{"mysql", "", "passwd", "server", "name"}, "", EMPTY_DB_USER_OR_PASSWD},
		{DBConfig{"mysql", "user", "", "server", "name"}, "", EMPTY_DB_USER_OR_PASSWD},
		{DBConfig{"mysql", "user", "passwd", "", "name"}, "", EMPTY_DB_SERVER},
		{DBConfig{"mysql", "user", "passwd", "server", ""}, "", EMPTY_DB_NAME},
	}

	for _, c := range cases {
		conn, err := c.cfg.GetConnection()
		assert.Equal(t, c.conn, conn, "Failed to get connection url: "+c.conn)
		assert.Equal(t, c.err, err)
	}
}

func TestLoadConfigFile(t *testing.T) {
	cases := []struct {
		name     string
		expected bool
	}{
		{"default.yml", true},
		{"non-exist.yml", false},
		{"invalid.yml", false},
	}

	for _, c := range cases {
		_, err := LoadConfigFile(filepath.Join("testdata", c.name))
		assert.Equal(t, c.expected, err == nil, "Failed to load config file: "+c.name)
	}
}
