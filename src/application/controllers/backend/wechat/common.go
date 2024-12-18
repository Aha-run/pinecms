package wechat

import (
	"errors"
	"github.com/xiusin/pine/contracts"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/xiusin/pine/di"
	"github.com/xiusin/pinecms/src/application/controllers"
	"github.com/xiusin/pinecms/src/application/models/tables"
	"github.com/xiusin/pinecms/src/common/helper"
)

func GetOfficialAccount(appid string) (*officialaccount.OfficialAccount, *tables.WechatAccount) {
	accountData := &tables.WechatAccount{}
	orm := helper.GetORM()
	orm.Where("app_id = ?", appid).Get(accountData)
	if accountData.Id == 0 {
		panic(errors.New("公众号" + appid + "不存在"))
	}
	wc, memory := wechat.NewWechat(), &WechatTokenCacher{Cache: di.MustGet(controllers.ServiceICache).(contracts.Cache)}
	cfg := &offConfig.Config{
		AppID:          accountData.AppId,
		AppSecret:      accountData.Secret,
		Token:          accountData.Token,
		EncodingAESKey: accountData.AesKey,
		Cache:          memory,
	}
	account := wc.GetOfficialAccount(cfg)
	return account, accountData
}

func SaveCacheMaterialListKey(key string, cacher contracts.Cache) {
	var keys []string
	cacher.GetWithUnmarshal(CacheKeyWechatMaterialListKeys, &keys)
	for _, cacheKey := range keys {
		if cacheKey == key {
			return
		}
	}
	keys = append(keys, key)
	cacher.SetWithMarshal(CacheKeyWechatMaterialListKeys, &keys, CacheTimeSecs)
}
