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
		reponame     string
		proto        string
		protoVersion string
		output       []string
		expected     bool
	}{
		{"notexist", "oci", "v1", nil, true},
		{"second/second", "oci", "v1", []string{"v0.1", "v0.2"}, true},
		{"second/second", "docker", "v1", nil, true},
		{"second/second", "oci", "v2", nil, true},
	}

	for _, c := range cases {
		tags, err := QueryTagsList(c.reponame, c.proto, c.protoVersion)
		assert.Equal(t, c.output, tags)
		assert.Equal(t, c.expected, err == nil)
	}
}

func TestQueryReposList(t *testing.T) {
	if !testReady {
		return
	}
	// init the testdata
	FreeTestDBData()
	InitTestDBData()

	testdata := []string{"first", "second/second", "third/third", "fourth/fourth/fourth"}

	repos, err := QueryReposList()
	assert.Equal(t, testdata, repos)
	assert.Nil(t, err)
}

func TestAddRepo(t *testing.T) {
	if !testReady {
		return
	}

	cases := []struct {
		reponame string
		expected bool
	}{
		{"TestAddRepo-A", true},
		// exist
		{"TestAddRepo-A", true},
	}
	for _, c := range cases {
		_, err := AddRepo(c.reponame)
		assert.Equal(t, c.expected, err == nil)
	}
}

func TestAddImage(t *testing.T) {
	if !testReady {
		return
	}

	cases := []struct {
		reponame     string
		tag          string
		proto        string
		protoVersion string
		expected     bool
	}{
		{"TestAddImage-A", "0.1", "test", "vtest", true},
		{"TestAddImage-A", "0.1", "test", "vtest", true},
	}
	for _, c := range cases {
		_, err := AddImage(c.reponame, c.tag, c.proto, c.protoVersion)
		assert.Equal(t, c.expected, err == nil)
	}
}
