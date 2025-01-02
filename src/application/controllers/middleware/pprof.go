package middleware

import (
	"github.com/valyala/fasthttp/pprofhandler"
	"github.com/xiusin/pine"
)

func Pprof(shouldUsePprofHandlerFn func(ctx *pine.Context) bool) pine.Handler {
	return func(ctx *pine.Context) {
		if shouldUsePprofHandlerFn != nil && shouldUsePprofHandlerFn(ctx) {
			pprofhandler.PprofHandler(ctx.RequestCtx)
			ctx.Stop()
			return
		}
		ctx.Next()
	}
}
