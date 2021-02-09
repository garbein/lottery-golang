package middleware

import (
	"net/http"

	"github.com/garbein/lottery-golang/responses"
	"github.com/garbein/lottery-golang/utils"
	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("AUTHORIZATION")
		if token == "" {
			c.JSON(http.StatusOK, responses.ErrorResponse("请先登录"))
			return
		}

		claims, err := utils.ParseAccessToken(token)

		if err != nil {
			c.JSON(http.StatusOK, responses.ErrorResponse("请先登录"))
			return
		}
		c.Set("claims", claims)
		c.Set("userId", claims.UserId)
		c.Next()
	}
}
