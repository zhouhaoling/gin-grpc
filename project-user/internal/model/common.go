package model

type BaseModel struct {
	Id  int64 `json:"id"`                    //id，主键
	MId int64 `json:"mid" gorm:"column:mid"` //用户id
}
