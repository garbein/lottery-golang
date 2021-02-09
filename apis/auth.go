package apis

import (
	"net/http"

	"github.com/garbein/lottery-golang/forms"
	"github.com/garbein/lottery-golang/responses"
	"github.com/garbein/lottery-golang/services"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {

	var loginForm forms.LoginForm

	if err := c.ShouldBindJSON(&loginForm); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(err.Error()))
		return
	}
	token, err := services.Login(loginForm.Mobile)
	if err != nil {
		c.JSON(http.StatusOK, responses.ErrorResponse("登录失败"))
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(gin.H{"token": token}))
	return
}
