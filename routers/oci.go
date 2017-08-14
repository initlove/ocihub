package router

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
)

const (
	ociPrefix = "/oci/v1"
)

func init() {
	RegisterRouter(ociPrefix, OCINameSpace())
}

func OCINameSpace() *beego.Namespace {
	ns := beego.NewNamespace(ociPrefix,
		beego.NSCond(func(ctx *context.Context) bool {
			logs.Debug("We get v1")
			return true
		}),
		beego.NSGet("/", func(ctx *context.Context) {
			ctx.Output.Body([]byte("ok"))
		}),
	)

	return ns
}
