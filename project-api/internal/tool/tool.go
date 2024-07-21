package tool

import "github.com/gin-gonic/gin"

func GetIP(c *gin.Context) string {
	ip := c.ClientIP()
	if ip == "::1" {
		return "127.0.0.1"
	}
	return ip
}
