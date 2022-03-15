package middleware

import (
	"fmt"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/xiusin/pine"
	"github.com/xiusin/pine/di"
	"github.com/xiusin/pinecms/src/application/controllers"
	"github.com/xiusin/pinecms/src/application/models/tables"
	"github.com/xiusin/pinecms/src/config"
	"strings"
	"xorm.io/xorm"
)

// VerifyJwtToken 验证jwt
// excludePrefix 需要排除认证的路由前缀
func VerifyJwtToken() pine.Handler {
	return func(ctx *pine.Context) {
		uri := ctx.Path()
		if !strings.Contains(uri, "login") &&
			!strings.Contains(uri, "captcha") &&
			!strings.Contains(uri, "thumb") {
			token := ctx.Header("Authorization")
			if token == "" {
				token, _ = ctx.GetString("token")
			}
			var hs = jwt.NewHS256([]byte(config.App().JwtKey))
			// 验证token
			var pl controllers.LoginAdminPayload
			_, err := jwt.Verify([]byte(token), hs, &pl)
			if err != nil {
				_ = ctx.Render().JSON(pine.H{"code": 1, "msg": "授权失败, 请重新登录"})
				return
			}

			ctx.SetUserValue("adminid", pl.AdminId)
			ctx.SetUserValue("roleid", pl.RoleID)

			if strings.Contains(uri, "user/info") {
				ctx.QueryArgs().Set("id", fmt.Sprintf("%d", pl.AdminId))
			}
			if !strings.Contains(uri, "/log/list") {
				di.MustGet(&xorm.Engine{}).(*xorm.Engine).Insert(&tables.RequestLog{
					Uri:      string(ctx.RequestURI()),
					Userid:   pl.AdminId,
					Params:   string(ctx.PostBody()),
					Username: pl.AdminName,
					Ip:       ctx.ClientIP(),
					Method:   string(ctx.Method()),
				})
			}
		}
		ctx.Next()
	}
}
