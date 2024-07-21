package model

type BaseModel struct {
	Id  int64 //id，主键
	MId int64 `gorm:"column:mid"` //用户id
}
