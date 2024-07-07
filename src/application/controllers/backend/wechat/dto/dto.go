package dto

type RuleSearch struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Param    struct {
		Appid string `json:"appid"`
	} `json:"param"`
}
