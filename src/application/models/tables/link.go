package tables

type Link struct {
	Id int64 `xorm:"int(11) autoincr not null pk 'id'" json:"id" schema:"id"`
	Linktype int `xorm:"tinyint(3) not null 'linktype'" json:"linktype" schema:"linktype" validate:"required"`
	Name string `xorm:"varchar(50) not null 'name'" json:"name" schema:"name" validate:"required"`
	Url string `xorm:"varchar(255) not null 'url'" json:"url" schema:"url" validate:"required"`
	Logo string `xorm:"varchar(100) not null 'logo'" json:"logo" schema:"logo" validate:"required"`
	Introduce string `xorm:"varchar(255) not null 'introduce'" json:"introduce" schema:"introduce" validate:"required"`
	Listorder int64 `xorm:"int(11) not null 'listorder'" json:"listorder" schema:"listorder" validate:"required"`
	Passed int `xorm:"tinyint(1) not null default '0' 'passed'" json:"passed" schema:"passed" validate:"required"`
	Addtime LocalTime `xorm:"datetime default 'null' 'addtime'" json:"addtime" schema:"addtime" validate:"required"`
}
