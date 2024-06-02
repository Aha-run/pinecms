package middleware

import (
	"github.com/arl/statsviz"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"github.com/xiusin/pine"
)

func StatesViz(app *pine.Application) pine.Handler {
	sev, _ := statsviz.NewServer()
	indexHandler := fasthttpadaptor.NewFastHTTPHandler(sev.Index())
	wsHandler := fasthttpadaptor.NewFastHTTPHandler(sev.Ws())

	app.GET("/debug/statsviz/*filepath", func(ctx *pine.Context) {
		if ctx.Params().Get("filepath") == "ws" {
			wsHandler(ctx.RequestCtx)
		} else {
			indexHandler(ctx.RequestCtx)
		}
	})
	return func(ctx *pine.Context) {
		ctx.Next()
	}
}
