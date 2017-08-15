package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"

	"github.com/initlove/ocihub/config"
	"github.com/initlove/ocihub/logger"
	"github.com/initlove/ocihub/models"
	"github.com/initlove/ocihub/routers"
	"github.com/initlove/ocihub/storage"
	_ "github.com/initlove/ocihub/storage/driver/filesystem"
)

func main() {
	if err := config.InitConfigFromFile("conf/ocihub.yml"); err != nil {
		return
	}
	if err := logger.InitLogger(); err != nil {
		logs.Warning(err)
	}

	cfg := config.GetConfig().DB
	conn, _ := cfg.GetConnection()
	if err := models.InitDB(conn, cfg.Driver, "default"); err != nil {
		logs.Critical("Error in init db: ", err)
		return
	}

	if err := storage.InitStorage(); err != nil {
		logs.Critical("Error in init storage: ", err)
		return
	}

	// demo of how to put content
	var ctx context.Context
	storage.Driver().PutContent(ctx, "/a/b/c", []byte("dliang"))

	nss := router.GetNamespaces()
	for name, ns := range nss {
		logs.Debug("Namespace '%s' is enabled", name)
		beego.AddNamespace(ns)
	}

	beego.Run()
}
