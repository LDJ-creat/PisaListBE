package middleware

import (
	"net/http"
	"strings"

	"github.com/PisaListBE/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未授权，请先登录",
			})
			c.Abort() //终止当前请求的处理流程，并立即返回响应给客户端
			return
		}

		// 去掉Bearer前缀
		token = strings.Replace(token, "Bearer ", "", 1)
		claims, err := jwt.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "token无效或已过期",
			})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID) //将声明中的用户ID存储到Gin的上下文中，这样后续的处理函数就可以通过上下文获取到用户ID。
		c.Next()                       //示中间件处理完毕，允许请求继续向下传递到下一个中间件或者最终的请求处理函数。
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// JWT 验证逻辑
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}
		// ... 其他验证逻辑
		c.Next()
	}
}
