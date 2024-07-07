package wechat

import (
	"errors"

	"github.com/xiusin/pinecms/src/application/controllers/backend"
	"github.com/xiusin/pinecms/src/application/controllers/backend/wechat/dto"
	"github.com/xiusin/pinecms/src/application/models/tables"
	"github.com/xiusin/pinecms/src/common/helper"
	"xorm.io/xorm"
)

type WechatRuleController struct {
	backend.BaseController
}

func (c *WechatRuleController) Construct() {
	c.Table = &tables.WechatMsgReplyRule{}
	c.Entries = &[]tables.WechatMsgReplyRule{}
	c.SearchFields = []backend.SearchFieldDsl{
		{Field: "appid"},
	}
	c.BaseController.Construct()
	c.OpBefore = c.before
}

func (c WechatRuleController) before(act int, params any) error {
	if act == backend.OpEdit || act == backend.OpAdd {
		//sess := params.(*xorm.Session).Clone()
		sess := params.(*xorm.Session)
		data := c.Table.(*tables.WechatMsgReplyRule)
		sess.Where("Match_Value = ?", data.MatchValue).Where("appid = ?", data.AppId)
		if act == backend.OpEdit {
			sess.Where("id = ?", data.Id)
		}
		if exist, _ := sess.Exist(&tables.WechatMsgReplyRule{}); exist {
			return errors.New("规则匹配值已经存在")
		}
	} else if act == backend.OpList {
		var search dto.RuleSearch
		helper.PanicErr(c.BindParse(&search))
		sess := params.(*xorm.Session)
		if search.Param.Appid != "" {
			sess.Where("appid = ?", search.Param.Appid)
		}
	}
	return nil
}
