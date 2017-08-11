package router

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
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
			fmt.Println("we get v1")
			return true
		}),
		beego.NSGet("/", func(ctx *context.Context) {
			ctx.Output.Body([]byte("ok"))
		}),
	)

	return ns
}
