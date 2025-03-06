package middleware

import (
	"Graduation-Project/common"
	"Graduation-Project/log"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请求未携带Token，无权限访问"})
			c.Abort()
			return
		}
		claims, err := ParseToken(tokenStr)
		if err != nil {
			log.Logger.Errorf(fmt.Sprintf("[%s]", err.Error()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": common.UNAUTHORIZED, "code": common.ERRCODE_UNAUTHORIZED})
			c.Abort()
			return
		}

		c.Set("jwtUsername", claims.Username)
		c.Set("id", claims.Id)
		c.Next()
	}
}
