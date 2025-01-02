package helper

import (
	"context"
	"log/slog"

	"github.com/xiusin/pine"
	"github.com/xiusin/pine/contracts"
	"github.com/xiusin/pine/di"
	"github.com/xiusin/pine/middlewares/traceid"
	"github.com/xiusin/pinecms/src/application/controllers"
)

// Inject 注入依赖
func Inject(key any, v any, single ...bool) {
	if len(single) == 0 {
		single = append(single, true)
	}
	if vi, ok := v.(di.BuildHandler); ok {
		di.Set(key, vi, single[0])
	} else {
		di.Set(key, func(_ di.AbstractBuilder) (i any, e error) {
			return v, nil
		}, single[0])
	}
}

// Cache 获取缓存服务
func Cache() contracts.Cache {
	return pine.Make(controllers.ServiceICache).(contracts.Cache)
}

// App 获取应用实例
func App() *pine.Application {
	return pine.Make(controllers.ServiceApplication).(*pine.Application)
}

// Slog 获取slog对象
func Slog(ctxs ...context.Context) *slog.Logger {
	logger := pine.Logger().(*slog.Logger)

	if len(ctxs) > 0 {
		requestID := ctxs[0].Value(traceid.Key)
		if requestID != nil {
			logger = logger.With(traceid.Key, requestID)
		}
	}
	return logger
}
