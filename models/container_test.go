package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryTagsList(t *testing.T) {
	if !testReady {
		return
	}

	cases := []struct {
		reponame      string
		proto         string
		proto_version string
		output        []string
		err           bool
	}{
		{"notexist", "oci", "v1", nil, true},
		{"second/second", "oci", "v1", []string{"v0.1", "v0.2"}, true},
		{"second/second", "docker", "v1", nil, true},
		{"second/second", "oci", "v2", nil, true},
	}

	for _, c := range cases {
		tags, err := QueryTagsList(c.reponame, c.proto, c.proto_version)
		assert.Equal(t, c.output, tags)
		assert.Equal(t, c.err, err == nil)
	}
}
