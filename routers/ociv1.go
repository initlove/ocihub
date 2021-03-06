package router

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"

	"github.com/initlove/ocihub/controllers"
)

const (
	ociV1Prefix = "/oci/v1"
)

func init() {
	if err := RegisterRouter(ociV1Prefix, OCIV1NameSpace()); err != nil {
		logs.Error("Failed to register router: '%s'.", ociV1Prefix)
	} else {
		logs.Debug("Register router '%s' registered.", ociV1Prefix)
	}
}

// OCIV1NameSpace defines the oci v1 router
func OCIV1NameSpace() *beego.Namespace {
	ns := beego.NewNamespace(ociV1Prefix,
		beego.NSCond(func(ctx *context.Context) bool {
			logs.Debug("We get ociv1")
			return true
		}),
		beego.NSGet("/", func(ctx *context.Context) {
			ctx.Output.Body([]byte("ok"))
		}),
		beego.NSGet("/_catalog", func(ctx *context.Context) {
			ctx.Output.Body([]byte("ok"))
		}),
		beego.NSRouter("/*/tags/list", &controllers.OCIV1Tag{}, "get:GetTagsList"),
		beego.NSRouter("/*/blobs/uploads/?:uuid", &controllers.OCIV1Blob{}, "post:PostBlob;patch:PatchBlob;put:PutBlob"),
		beego.NSRouter("/*/blobs/:digest", &controllers.OCIV1Blob{}, "head:HeadBlob;get:GetBlob;delete:DeleteBlob"),
		beego.NSRouter("/*/manifest/:tags", &controllers.OCIV1Manifest{}, "get:GetManifest;put:PutManifest;delete:DeleteManifest"),
	)

	return ns
}
