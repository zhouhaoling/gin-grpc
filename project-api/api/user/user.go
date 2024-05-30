package user

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"test.com/common"
)

type HandlerUser struct {
}

func NewHandlerUser() *HandlerUser {
	return &HandlerUser{}
}

func (u *HandlerUser) getCaptcha(c *gin.Context) {
	result := common.NewResponseData()
	mobile := c.PostForm("mobile")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	LoginServiceClient.GetCaptcha(ctx)
}
