package wechat

import (
	"time"

	cache2 "github.com/xiusin/pine/cache"
)

type WechatTokenCacher struct {
	cache2.AbstractCache
}

func (w WechatTokenCacher) Get(key string) any {
	byts, err := w.AbstractCache.Get(key)
	if err != nil {
		return nil
	}
	return string(byts)
}

func (w WechatTokenCacher) Set(key string, val any, timeout time.Duration) error {
	return w.AbstractCache.Set(key, []byte(val.(string)), int(timeout.Seconds()))
}

func (w WechatTokenCacher) IsExist(key string) bool {
	return w.AbstractCache.Exists(key)
}

func (w WechatTokenCacher) Delete(key string) error {
	return w.AbstractCache.Delete(key)
}
