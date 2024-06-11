package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/casbin/casbin/v2"
	xd "github.com/casbin/xorm-adapter"
	"xorm.io/xorm"

	"github.com/xiusin/pine"
	"github.com/xiusin/pine/di"
	"github.com/xiusin/pinecms/src/application/controllers"
	"github.com/xiusin/pinecms/src/application/models/tables"
	"github.com/xiusin/pinecms/src/common/helper"
)

func Casbin(engine *xorm.Engine, conf string) pine.Handler {
	var _locker = &sync.Mutex{}
	adapter, err := xd.NewAdapterByEngine(engine)
	helper.PanicErr(err)
	enforcer, err := casbin.NewEnforcer(helper.GetRootPath(conf), adapter)
	helper.PanicErr(err)
	di.Set(controllers.ServiceCasbinEnforcer, func(builder di.AbstractBuilder) (any, error) {
		return enforcer, nil
	}, true)
	di.Set(controllers.ServiceCasbinClearPolicy, func(builder di.AbstractBuilder) (any, error) {
		return clearPolicy(enforcer, _locker), nil
	}, true)

	addPolicyHandler := addPolicy(engine, enforcer, _locker)
	addPolicyHandler()

	di.Instance(controllers.ServiceCasbinAddPolicy, addPolicyHandler)

	return func(ctx *pine.Context) {
		adminId := ctx.Value("adminid")
		if adminId != nil {
			var admin = &tables.Admin{}
			exist, _ := engine.Where("id = ?", adminId).Get(admin)
			pathString := strings.Split(strings.Trim(ctx.Path(), "/"), "/")

			if exist && len(pathString) >= 3 && pathString[0] == "v2" {
				roles := admin.RoleIdList
				var passable bool
				for _, role := range roles {
					passable, _ = enforcer.Enforce(fmt.Sprintf("%d", role), pathString[1], pathString[2])
					if passable {
						ctx.Next()
						return
					}
				}
			}
			ctx.Abort(http.StatusForbidden)
			return
		}
		ctx.Next()
	}
}

func clearPolicy(enforcer *casbin.Enforcer, _locker *sync.Mutex) func() {
	return func() {
		_locker.Lock()
		defer _locker.Unlock()
		enforcer.ClearPolicy()
	}
}

// 根据角色注入权限
func addPolicy(engine *xorm.Engine, enforcer *casbin.Enforcer, _locker *sync.Mutex) func() {
	return func() {
		_locker.Lock()
		defer _locker.Unlock()
		_, _ = engine.Table(&xd.CasbinRule{}).Where("v0 > 0").Delete()
		var roles []tables.AdminRole
		_ = engine.Find(&roles)
		var menus []*tables.Menu
		_ = engine.Find(&menus)

		var menuMap = map[int64]*tables.Menu{}
		for _, menu := range menus {
			menuMap[menu.Id] = menu
		}
		for _, role := range roles {
			var privs []tables.AdminRolePriv
			engine.Where("role_id = ?", role.Id).Find(&privs)
			for _, priv := range privs {
				menu := menuMap[priv.MenuId]
				if menu == nil {
					continue
				}

				args := []string{fmt.Sprintf("%d", role.Id)}
				args = append(args, strings.Split(strings.Trim(menu.Router, "/"), "/")...)
				enforcer.AddPolicy(helper.ConvertToAnySlice(args)...)
			}
		}
		enforcer.SavePolicy()
	}
}
