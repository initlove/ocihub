package filesystem

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/astaxie/beego/context"
	"github.com/stretchr/testify/assert"
)

func TestPut(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "ocidata")
	defer os.RemoveAll(tmpDir)

	paras := make(map[string]interface{})
	paras["rootDirectory"] = tmpDir

	var ctx context.Context
	var d driver
	d.Init(paras)

	testPath := "testPath"
	testData := []byte("testdata")

	d.PutContent(ctx, testPath, testData)
	data, _ := d.GetContent(ctx, testPath)
	assert.Equal(t, testData, data)
}
