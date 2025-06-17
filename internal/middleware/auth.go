package middleware

import (
	"go-gin/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Token")
		if token != "Token123" { // 这里的Token123需要换为具体生成的token
			utils.GetResponseJson(c, utils.TOKEN_ERROR, "权限不足", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
