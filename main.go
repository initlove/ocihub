package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/initlove/ocihub/config"
	"github.com/initlove/ocihub/logger"
	"github.com/initlove/ocihub/models"
	"github.com/initlove/ocihub/routers"
	"github.com/initlove/ocihub/session"
	"github.com/initlove/ocihub/storage"

	_ "github.com/initlove/ocihub/session/memory"
	_ "github.com/initlove/ocihub/storage/driver/filesystem"
)

func main() {
	cfg, err := config.InitConfigFromFile("conf/ocihub.yml")
	if err != nil {
		return
	}
	if err := logger.InitLogger(cfg.Log); err != nil {
		logs.Warning(err)
	}

	conn, _ := cfg.DB.GetConnection()
	if err := models.InitDB(conn, cfg.DB.Driver, "default"); err != nil {
		logs.Critical("Error in init db: ", err)
		return
	}

	if err := storage.InitStorage(cfg.Storage); err != nil {
		logs.Critical("Error in init storage: ", err)
		return
	}

	if err := session.InitSession(cfg.Session); err != nil {
		logs.Critical("Error in init session: ", err)
		return
	}

	nss := router.GetNamespaces()
	for name, ns := range nss {
		logs.Debug("Namespace '%s' is enabled", name)
		beego.AddNamespace(ns)
	}

	beego.Run()
}
