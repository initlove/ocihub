package models

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var testReady = false

func init() {
	err := InitTestDB()
	if err != nil {
		fmt.Println("Fail to test sql,skip all the sql testing: ", err)
		return
	}

	// FIXME: we need a better testing framework to free data after everything is done.
	// Free anyway before init the data
	FreeTestDBData()
	err = InitTestDBData()
	if err != nil {
		fmt.Println("Fail to init data, mostly you need to check your test sql: ", err)
		return
	}

	testReady = true
}

// InitTestDB
//  conn should be something like root:1234@tcp(localhost:3306)/test?charset=utf8
func InitTestDB() error {
	conn := os.Getenv("TESTCONN")
	if conn == "" {
		return errors.New("Please set the 'TESTCONN' before using db unit test")
	}

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.DefaultTimeLoc = time.UTC

	err := orm.RegisterDataBase("default", "mysql", conn)
	if err == nil {
		orm.RunSyncdb("default", false, true)
		return nil
	}

	return errors.New("Cannot connect to the test database")
}

func InitTestDBData() error {
	initsql, err := ioutil.ReadFile("testdata/init.sql")
	if err != nil {
		return fmt.Errorf("Failed to init test data: %v", err)
	}

	// Seems orm can only exec one sql at one time.
	for _, sql := range strings.SplitAfter(string(initsql), ";") {
		if len(strings.TrimSpace(sql)) == 0 {
			continue
		}
		_, err = orm.NewOrm().Raw(sql).Exec()
		if err != nil {
			break
		}
	}
	if err != nil {
		FreeTestDBData()
		return fmt.Errorf("Failed to exec init sql: %v", err)
	}

	return nil
}

func FreeTestDBData() error {
	freesql, err := ioutil.ReadFile("testdata/free.sql")
	if err != nil {
		return err
	}
	for _, sql := range strings.SplitAfter(string(freesql), ";") {
		if len(strings.TrimSpace(sql)) == 0 {
			continue
		}
		orm.NewOrm().Raw(sql).Exec()
	}

	return nil
}

func TestInitDB(t *testing.T) {
	if !testReady {
		return
	}

	cases := []struct {
		conn     string
		driver   string
		name     string
		expected bool
	}{
		{os.Getenv("TESTCONN"), "mysql", "TestInitDB-case-0", true},
		{os.Getenv("TESTCONN"), "liangsql", "TestInitDB-case-1", false},
		{"localhost:22", "mysql", "TestInitDB-case-2", false},
	}

	for _, c := range cases {
		err := InitDB(c.conn, c.driver, c.name)
		assert.Equal(t, c.expected, err == nil)
	}
}
