package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"

	"github.com/initlove/ocihub/config"
)

func InitDB() error {
	var err error
	cfg := config.GetConfig().DB
	conn, err := cfg.GetConnection()
	if err != nil {
		return err
	}

	if cfg.Driver != "mysql" {
		return errors.New("Only support mysql yet.")
	}

	orm.RegisterDriver(cfg.Driver, orm.DRMySQL)
	orm.DefaultTimeLoc = time.UTC

	for i := 0; i < 3; i++ {
		<-time.After(time.Second * 5)
		err = orm.RegisterDataBase("default", cfg.Driver, conn)
		if err == nil {
			orm.SetMaxIdleConns("default", 30)
			orm.SetMaxOpenConns("default", 30)
			orm.RunSyncdb("default", false, true)
			return nil
		}
		logs.Debug("Try to register database for %d times...", i)
	}

	return errors.New("Cannot connect to database")
}
