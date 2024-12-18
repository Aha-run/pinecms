package cmd

import (
	"errors"
	"fmt"
	"github.com/xiusin/pine"

	. "github.com/xiusin/pinecms/src/config"

	"github.com/spf13/cobra"
	"github.com/xiusin/pinecms/src/application/models/tables"
	"xorm.io/xorm"
)

var menuCmd = &cobra.Command{
	Use:   "menu",
	Short: "生成模块菜单权限",
	Run: func(cmd *cobra.Command, args []string) {
		InitDB() // 方法不可放到init里，否则缓存组件阻塞
		if !IsDebug() {
			pine.Logger().Info("非Debug模式，不支持 Menu 命令")
			return
		}
		table, _ := cmd.Flags().GetString("table")
		force, _ := cmd.Flags().GetBool("force")
		menu, _ := cmd.Flags().GetString("menu")
		if table == "" {
			cmd.Help()
			return
		}
		if menu == "" {
			menu = table + "管理"
		}

		menus := []struct {
			MenuName string
			A        string
		}{
			{
				MenuName: menu,
				A:        "list",
			},

			{
				MenuName: "查看",
				A:        "info",
			},

			{
				MenuName: "新增",
				A:        "add",
			},

			{
				MenuName: "编辑",
				A:        "edit",
			},

			{
				MenuName: "删除",
				A:        "delete",
			},
		}
		role := &tables.Menu{}
		if !force {
			count, _ := Orm().Table(role).Where("c = ?", table).Count()
			if count > 0 {
				pine.Logger().Error(fmt.Sprintf("已经存在%s的相关菜单, 如需强制覆盖请追加参数--force true", table))
				return
			}
		}

		Orm().Where("c = ?", table).Delete(role)

		_, err := Orm().Transaction(func(session *xorm.Session) (any, error) {
			var parId int64
			for k, item := range menus {
				role := tables.Menu{}
				role.Id = 0
				role.Name = item.MenuName
				role.Identification = table + ":" + item.A
				if k == 0 {
					role.Parentid = 0
					role.Display = true
				} else {
					role.Parentid = parId
					role.Display = false
				}
				id, err := Orm().Insert(&role)
				if err != nil {
					return nil, err
				}
				if role.Id == 0 {
					fmt.Printf("%d, %#v", id, role)
					return nil, errors.New("生成菜单失败")
				}
				if k == 0 {
					parId = role.Id
				}
			}
			return nil, nil
		})

		if err != nil {
			pine.Logger().Error(err.Error())
		}
	},
}

func init() {
	menuCmd.Flags().String("table", "", "数据库表名")
	menuCmd.Flags().String("menu", "", "顶级菜单名称")
	menuCmd.Flags().Bool("force", false, "是否强制生成")
}
