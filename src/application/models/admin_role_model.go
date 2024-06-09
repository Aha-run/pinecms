package models

import (
	"log"

	"xorm.io/xorm"

	"github.com/xiusin/pinecms/src/application/models/tables"
	"github.com/xiusin/pinecms/src/common/helper"
)

type AdminRoleModel struct {
	orm *xorm.Engine
}

func NewAdminRoleModel() *AdminRoleModel {
	return &AdminRoleModel{orm: helper.GetORM()}
}

func (a *AdminRoleModel) List(page, rows int) ([]tables.AdminRole, int64) {
	start := (page - 1) * rows
	roles := []tables.AdminRole{}
	total, _ := a.orm.Limit(rows, start).FindAndCount(&roles)
	return roles, total
}

func (a *AdminRoleModel) All() map[int64]*tables.AdminRole {
	var roles = map[int64]*tables.AdminRole{}
	var _roles []*tables.AdminRole
	a.orm.Find(&_roles)
	for _, role := range _roles {
		roles[role.Id] = role
	}
	return roles
}

func (a *AdminRoleModel) CheckRoleName(id int64, rolename string) bool {
	exists, _ := a.orm.Where("roleid <> ?", id).Where("rolename = ?", rolename).Exist()
	return exists
}

func (a *AdminRoleModel) AddRole(rolename, description string, disabled, listorder int64) bool {
	newRole := tables.AdminRole{
		Rolename:    rolename,
		Description: description,
		Disabled:    disabled,
		Listorder:   listorder,
	}
	insertId, err := a.orm.Insert(&newRole)
	if err != nil || insertId == 0 {
		log.Println(err, insertId)
		return false
	}
	return true
}

func (a *AdminRoleModel) GetRoleById(id int64) (tables.AdminRole, error) {
	role := tables.AdminRole{Id: id}
	_, err := a.orm.Get(&role)
	return role, err
}

func (a *AdminRoleModel) UpdateRole(role tables.AdminRole) bool {
	res, err := a.orm.Where("id=?", role.Id).MustCols("disabled").Update(&role)
	if err != nil || res == 0 {
		return false
	}
	return true
}

func (a *AdminRoleModel) DeleteRole(role tables.AdminRole) bool {
	res, err := a.orm.Delete(&role)
	if err != nil || res == 0 {
		return false
	}
	return true
}
