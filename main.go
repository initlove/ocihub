package main

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"

	"github.com/initlove/ocihub/config"
	"github.com/initlove/ocihub/models"
	"github.com/initlove/ocihub/routers"
	"github.com/initlove/ocihub/storage"
	_ "github.com/initlove/ocihub/storage/driver/filesystem"
)

func main() {
	cfg, err := config.LoadConfigFile("conf/ocihub.yml")
	if err != nil {
		return
	}
	if err := models.InitDB(cfg.DB); err != nil {
		fmt.Println("Error in init db: ", err)
		return
	}

	if err := storage.InitStorage(); err != nil {
		fmt.Println("Error in init storage: ", err)
		return
	}

	// demo of how to put content
	var ctx context.Context
	err = storage.Driver().PutContent(ctx, "/a/b/c", []byte("dliang"))

	nss := router.GetNamespaces()
	for name, ns := range nss {
		fmt.Println("we are adding ns: ", name)
		beego.AddNamespace(ns)
	}

	beego.Run()
}
