package main

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"

	"github.com/initlove/ocihub/config"
	"github.com/initlove/ocihub/models"
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

	ns := beego.NewNamespace("/v1",
		beego.NSCond(func(ctx *context.Context) bool {
			fmt.Println("we get v1")
			return true
		}),
		beego.NSGet("/version", func(ctx *context.Context) {
			fmt.Println("we get version")
			ctx.Output.Body([]byte("1.0"))
		}),
	)

	beego.AddNamespace(ns)
	beego.Run()
}
