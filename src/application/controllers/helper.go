package controllers

import (
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/xiusin/pine/di"
)

// LoginAdminPayload 登录JWT载体
type LoginAdminPayload struct {
	jwt.Payload
	Id        int64   `json:"id"`
	AdminId   int64   `json:"admin_id"`
	AdminName string  `json:"admin_name"`
	RoleID    []int64 `json:"role_id"`
}

// GetTableName 获取表名
func GetTableName(name string) string {
	tablePrefix := di.MustGet(ServiceTablePrefix).(string)
	return tablePrefix + name
}
