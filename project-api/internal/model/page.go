package model

import "github.com/gin-gonic/gin"

type PageStruct struct {
	Page     int64 `json:"page" form:"page"`
	PageSize int64 `json:"pageSize" form:"pageSize"`
}

func (p *PageStruct) Bind(c *gin.Context) {
	c.ShouldBind(&p)
	if p.Page == 0 {
		p.Page = 1
	}
	if p.PageSize == 0 {
		p.PageSize = 10
	}
	return
}
