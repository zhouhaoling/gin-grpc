package model

// Organization 数据库实体类，对应ms_organization表
type Organization struct {
	Id          int64
	Name        string
	Avatar      string
	Description string
	MemberId    int64
	CreateTime  int64
	Personal    int32
	Address     string
	Province    int32
	City        int32
	Area        int32
}

func (*Organization) TableName() string {
	return "ms_organization"
}
