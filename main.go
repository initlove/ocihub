package main

import (
	"fmt"

	"github.com/astaxie/beego"

	"github.com/initlove/ocihub/config"
	"github.com/initlove/ocihub/models"
	"github.com/initlove/ocihub/routers"
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

	nss := router.GetNamespaces()
	for name, ns := range nss {
		fmt.Println("we are adding ns: ", name)
		beego.AddNamespace(ns)
	}

	beego.Run()
}
