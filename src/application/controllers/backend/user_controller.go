package backend

import (
	"errors"
	"strings"

	"xorm.io/builder"
	"xorm.io/xorm"

	"github.com/xiusin/pinecms/src/application/models"
	"github.com/xiusin/pinecms/src/application/models/tables"
	"github.com/xiusin/pinecms/src/common/helper"
)

type UserController struct {
	BaseController
}

func (c *UserController) Construct() {
	c.KeywordsSearch = []SearchFieldDsl{
		{Field: "username", Op: "LIKE", DataExp: "%$?%"},
		{Field: "email", Op: "LIKE", DataExp: "%$?%"},
		{Field: "realname", Op: "LIKE", DataExp: "%$?%"},
	}

	c.SearchFields = []SearchFieldDsl{
		{Field: "departmentIds", CallBack: func(session *xorm.Session, i ...any) {
			session.Where(builder.In("department_id", i))
		}},
	}

	c.Orm = helper.GetORM()
	c.Table = &tables.Admin{}
	c.Entries = &[]*tables.Admin{}

	c.TableKey = "id"
	c.TableStructKey = "Userid"
	c.ExceptCols = []string{"password", "encrypt"}
	c.OpBefore = c.before
	c.OpAfter = c.after
}

func (c *UserController) GetAdminInfo() {
	c.GetInfo()
}
func (c *UserController) GetInfo() {
	c.Orm.Where("id = ?", c.Ctx().Value("adminid")).Get(c.Table)
	c.Table.(*tables.Admin).Password = ""
	helper.Ajax(c.Table, 0, c.Ctx())
}

func (c *UserController) PostPersonUpdate() {
	helper.PanicErr(c.Ctx().BindJSON(c.Table))
	user := c.Table.(*tables.Admin)
	loginUser := &tables.Admin{}
	c.Orm.Where("id = ?", c.Ctx().Value("adminid")).Get(loginUser)

	if loginUser.Userid == 0 {
		helper.Ajax("获取信息失败", 1, c.Ctx())
		return
	}

	loginUser.Avatar = user.Avatar
	loginUser.Realname = user.Realname
	if len(user.Password) > 0 {
		loginUser.Encrypt = helper.GetRandomString(6)
		loginUser.Password = helper.Password(user.Password, loginUser.Encrypt)
	}
	c.Orm.Where("id = ?", loginUser.Userid).Update(loginUser)
	helper.Ajax("更新信息成功", 0, c.Ctx())
}

func (c *UserController) before(opType int, param any) error {
	if opType == OpAdd || opType == OpEdit {
		p, exist := param.(*tables.Admin), false
		if p.Userid > 0 {
			exist, _ = c.Orm.Table(p).Where("id <> ? and (username = ? or email = ?)", p.Userid, p.Username, p.Email).Exist()
		} else {
			exist, _ = c.Orm.Table(p).Where("username = ? or email = ?", p.Username, p.Email).Exist()
		}
		if exist {
			return errors.New("用户名或邮箱已存在")
		}

		if len(p.Password) > 0 {
			p.Encrypt = helper.GetRandomString(6)
			p.Password = helper.Password(p.Password, p.Encrypt)
		}
	}
	if opType == OpDel { // 删除权限控制
		p := param.(*idParams)
		for _, id := range p.Ids {
			if 1 == id {
				return errors.New("无法删除内置超级管理员")
			}
		}
	}
	return nil
}

func (c *UserController) after(opType int, param any) error {
	if opType == OpList {
		admins := c.Entries.(*[]*tables.Admin)
		roles := models.NewAdminRoleModel().All()
		for _, admin := range *admins {
			for _, roleId := range admin.RoleIdList {
				if r, ok := roles[roleId]; ok {
					admin.RoleName += r.Rolename + ","
				}
			}
			admin.RoleName = strings.TrimRight(admin.RoleName, ",")
		}
	}
	return nil
}

func (c *UserController) PostLogout() {
	helper.Ajax("退出成功", 0, c.Ctx())
}
