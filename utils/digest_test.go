package utils

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDigest(t *testing.T) {
	expected := "e637dbddc5be31a623306e9e21e4d7c77878ba425445649d2a6804755976f29d"
	data, _ := ioutil.ReadFile("testdata/" + expected)
	assert.Equal(t, expected, GetDigest("sha256", data))
}
