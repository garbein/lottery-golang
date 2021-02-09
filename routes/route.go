package routes

import (
	"net/http"

	"github.com/garbein/lottery-golang/apis"
	"github.com/garbein/lottery-golang/middleware"
	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
)

func InitRoute() *gin.Engine {

	r := gin.Default()

	r.Use(favicon.New("./favicon.ico"))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "welcome")
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/login", apis.Login)

	authorized := r.Group("/")

	authorized.Use(middleware.JWT())
	{
		authorized.POST("lottery", apis.StartLottery)
	}

	return r
}
