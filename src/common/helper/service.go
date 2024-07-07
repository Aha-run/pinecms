package helper

import (
	"github.com/xiusin/pine"
	"github.com/xiusin/pine/cache"
	"github.com/xiusin/pine/di"
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
		di.Set(key, func(builder di.AbstractBuilder) (i any, e error) {
			return v, nil
		}, single[0])
	}
}

// 获取缓存服务
func AbstractCache() cache.AbstractCache {
	return pine.Make(controllers.ServiceICache).(cache.AbstractCache)
}

// 获取应用实例
func App() *pine.Application {
	return pine.Make(controllers.ServiceApplication).(*pine.Application)
}
