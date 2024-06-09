package tables

type Advert struct {
	Id        int64      `json:"id"`
	Name      string     `json:"name" schema:"name"`
	SpaceID   int64      `xorm:"space_id" json:"space_id" schema:"space_id"`
	SpaceName string     `xorm:"-" json:"space_name"  schema:"-"`
	LinkUrl   string     `json:"link_url" schema:"link_url"`
	Image     string     `json:"image" schema:"image"`
	ListOrder uint       `xorm:"listorder default 0" json:"listorder" schema:"listorder"`
	StartTime *LocalTime `xorm:"datetime 'start_time' default NULL " json:"startTime"`
	EndTime   *LocalTime `xorm:"datetime 'end_time' default NULL " json:"endTime"`
	DateRange []string   `json:"date_range" schema:"-"`  // 展示日期范围
	Status    bool       `json:"status" schema:"status"` // 状态
}
