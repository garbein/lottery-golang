package apis

import (
	"net/http"

	"github.com/garbein/lottery-golang/responses"
	"github.com/garbein/lottery-golang/services"
	"github.com/gin-gonic/gin"
)

func StartLottery(c *gin.Context) {
	userId := c.GetInt("userId")
	if userId <= 0 {
		c.JSON(http.StatusOK, responses.ErrorResponse("请重新登录"))
		return
	}
	responseBody, err := services.StartLottery(userId)
	if err != nil {
		c.JSON(http.StatusOK, responses.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, responseBody)
	return
}
