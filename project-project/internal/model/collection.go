package model

type Collection struct {
	Id          int64
	ProjectCode int64
	MemberCode  int64
	CreateTime  int64
}

func (c *Collection) TableName() string {
	return "ms_project_collection"
}
