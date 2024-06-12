package project

import "github.com/gin-gonic/gin"

type HandlerProject struct {
}

func NewHandlerProject() *HandlerProject {
	return &HandlerProject{}
}

func (p *HandlerProject) index(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
	})
}
