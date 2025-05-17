package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/azan-boss/posty/internal/handler/auth"
	"github.com/azan-boss/posty/internal/utils/response"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token==""{
			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("unauthorized")))
			return
		}
		token=strings.TrimPrefix(token, "Bearer ")
		// fmt.Println(token)
		claims, err := auth.VerifyJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("unauthorized")))
			return
		}
		fmt.Println("claims",claims)
		c.Set("username", claims.Username)
		c.Set("userId", claims.UserId)
		c.Next()
	}
}
