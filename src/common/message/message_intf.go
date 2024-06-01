package message

const (
	TypeDefault = iota
	TypeRegister
	TypeLogin
	TypeModifyProfile
	TypeFindPwd
	TypeNotice
)

// AbstractMessage 发送接口
type AbstractMessage interface {
	Init() error
	Send(receiver []string, subject string, body string) error
	// receiver 接收人数组 params 模板内数据
	Notice(receiver []string, params []any, templateId int) error
	// 更新单例配置
	UpdateCfg() error
}

var MessageServiceDict = map[string]string{
	ServiceNullMessage: "空短信",
	//ServiceSmsMessage:   "阿里信息",
	ServiceEmailMessage: "邮箱",
}
