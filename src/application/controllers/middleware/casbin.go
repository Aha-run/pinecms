package middleware

import (
	_ "embed"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	xd "github.com/casbin/xorm-adapter"
	"xorm.io/xorm"

	"github.com/xiusin/pine"
	"github.com/xiusin/pine/di"
	"github.com/xiusin/pinecms/src/application/controllers"
	"github.com/xiusin/pinecms/src/application/models/tables"
	"github.com/xiusin/pinecms/src/common/helper"
)

var publicApis = []string{
	"/public/menu",
	"/user/admin_info",
	"/user/login",
	"/user/logout",
}

//go:embed rbac_models.conf
var _model string

func Casbin(engine *xorm.Engine) pine.Handler {
	var _locker = &sync.Mutex{}
	adapter, err := xd.NewAdapterByEngine(engine)
	helper.PanicErr(err)

	model, err := casbinModel.NewModelFromString(_model)
	helper.PanicErr(err)

	enforcer, err := casbin.NewSyncedCachedEnforcer(model, adapter)
	helper.PanicErr(err)
	enforcer.EnableCache(true)
	enforcer.EnableAutoSave(true)
	enforcer.StartAutoLoadPolicy(time.Minute * 2)

	helper.PanicErr(enforcer.LoadPolicy())

	helper.Inject(controllers.ServiceCasbinEnforcer, enforcer, true)
	di.Set(controllers.ServiceCasbinClearPolicy, func(_ di.AbstractBuilder) (any, error) {
		return clearPolicy(enforcer, _locker), nil
	}, true)

	addPolicyHandler := addPolicy(engine, enforcer, _locker)
	addPolicyHandler()
	di.Instance(controllers.ServiceCasbinAddPolicy, addPolicyHandler)

	var paMaps = map[string]struct{}{}
	for _, v := range publicApis {
		paMaps[v] = struct{}{}
	}

	return func(ctx *pine.Context) {
		adminId := ctx.Value("adminid")
		if adminId != nil {
			var admin = &tables.Admin{}
			if exist, _ := engine.Where("id = ?", adminId).Get(admin); !exist {
				ctx.Abort(http.StatusForbidden)
			}
			pathString := strings.Split(strings.Trim(ctx.Path(), "/"), "/")

			if len(pathString) >= 3 && pathString[0] == "v2" {
				pathString[0] = ""
				if _, ok := paMaps[strings.Join(pathString, "/")]; ok {
					ctx.Next()
					return
				}
				for _, role := range admin.RoleIdList {
					if passable, _ := enforcer.Enforce(fmt.Sprintf("%d", role), pathString[1], pathString[2]); passable {
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

func clearPolicy(enforcer *casbin.SyncedCachedEnforcer, _locker *sync.Mutex) func() {
	return func() {
		_locker.Lock()
		defer _locker.Unlock()
		enforcer.ClearPolicy()
	}
}

// 根据角色注入权限
func addPolicy(engine *xorm.Engine, enforcer *casbin.SyncedCachedEnforcer, _locker *sync.Mutex) func() {
	return func() {
		_locker.Lock()
		defer _locker.Unlock()
		engine.Table(&xd.CasbinRule{}).Where("v0 > 0").Delete()
		var roles []tables.AdminRole
		engine.Find(&roles)
		var menus []tables.Menu
		engine.Find(&menus)

		var menuMap = map[int64]tables.Menu{}
		for _, menu := range menus {
			menuMap[menu.Id] = menu
		}
		for _, role := range roles {
			for _, menuId := range role.MenuIdList {
				if menu, ok := menuMap[menuId]; ok {
					args := []string{fmt.Sprintf("%d", role.Id)}
					args = append(args, strings.Split(strings.Trim(menu.Router, "/"), "/")...)
					enforcer.AddPolicy(helper.ConvertToAnySlice(args)...)
				}
			}
		}
		enforcer.SavePolicy()
	}
}
