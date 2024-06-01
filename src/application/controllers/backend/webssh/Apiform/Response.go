package Apiform

type Resp struct {
	Code  int    `json:"code"`
	Data  any    `json:"data"`
	Msg   string `json:"msg"`
	Token string `json:"token"`
}
