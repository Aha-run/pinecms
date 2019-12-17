package tables

type IriscmsAdmin struct {
	Userid        int64  `xorm:"pk"`
	Username      string `xorm:"unique"`
	Password      string
	Roleid        int64
	Encrypt       string
	Lastloginip   string
	Lastlogintime int64
	Email         string
	Realname      string
}
