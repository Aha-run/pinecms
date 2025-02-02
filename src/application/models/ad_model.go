package models

import (
	"github.com/xiusin/pine"
	"github.com/xiusin/pinecms/src/application/models/tables"
	"github.com/xiusin/pinecms/src/common/helper"
	"xorm.io/xorm"
)

type AdModel struct {
	orm *xorm.Engine
}

func NewAdModel() *AdModel {
	return &AdModel{orm: helper.GetORM()}
}

func (l *AdModel) GetList(page, limit int) ([]tables.Advert, int64) {
	offset := (page - 1) * limit
	var list = []tables.Advert{}
	var total int64
	var err error
	if total, err = l.orm.Desc("listorder").Limit(limit, offset).FindAndCount(&list); err != nil {
		pine.Logger().Error(err.Error())
	}
	return list, total
}

func (l *AdModel) Add(data *tables.Advert) int64 {
	id, err := l.orm.InsertOne(data)
	if err != nil {
		pine.Logger().Error(err.Error())
	}
	return id
}

func (l *AdModel) Delete(id []int64) bool {
	res, err := l.orm.ID(id).Delete(&tables.Advert{})
	if err != nil {
		pine.Logger().Error(err.Error())
	}
	return res > 0
}

func (l *AdModel) Get(id int64) *tables.Advert {
	var link = &tables.Advert{}
	ok, _ := l.orm.ID(id).Get(link)
	if !ok {
		return nil
	}
	return link
}

func (l *AdModel) Update(data *tables.Advert) bool {
	id, err := l.orm.ID(data.Id).Update(data)
	if err != nil {
		pine.Logger().Error(err.Error())
	}

	return id > 0
}
