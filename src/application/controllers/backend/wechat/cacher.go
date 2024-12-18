package wechat

import (
	"github.com/xiusin/pine/contracts"
	"time"
)

type WechatTokenCacher struct {
	contracts.Cache
}

func (w WechatTokenCacher) Get(key string) any {
	byts, err := w.Cache.Get(key)
	if err != nil {
		return nil
	}
	return string(byts)
}

func (w WechatTokenCacher) Set(key string, val any, timeout time.Duration) error {
	return w.Cache.Set(key, []byte(val.(string)), int(timeout.Seconds()))
}

func (w WechatTokenCacher) IsExist(key string) bool {
	return w.Cache.Exists(key)
}

func (w WechatTokenCacher) Delete(key string) error {
	return w.Cache.Delete(key)
}
