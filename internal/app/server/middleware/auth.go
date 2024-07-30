package middleware

import (
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// todo : 登录用户AUTH校验

		c.Set("uid", 1)
		c.Set("username", "jiangyu")
		c.Set("name", "江屿")
		c.Next()
	}
}