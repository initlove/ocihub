package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
)

//TODO: more logs info

func CTX_ERROR_WRAP(ctx *context.Context, code int, err error, msg string) {
	ctx.Output.SetStatus(code)
	ctx.Output.Body([]byte(msg))

	if err != nil {
		logs.Trace("Failed to [%s] [%s] [%d]: [%v]", ctx.Input.Method(), ctx.Input.URI(), code, err)
	} else {
		logs.Trace("Failed to [%s] [%s] [%d]", ctx.Input.Method(), ctx.Input.URI(), code)
	}
}

func CTX_SUCCESS_WRAP(ctx *context.Context, code int, result interface{}, header map[string]string) {
	ctx.Output.SetStatus(code)
	for n, v := range header {
		ctx.Output.Header(n, v)
	}
	output, _ := json.Marshal(result)
	ctx.Output.Body(output)

	logs.Trace("Succeed in [%s] [%s].", ctx.Input.Method(), ctx.Input.URI())
}
