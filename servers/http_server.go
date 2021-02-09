package servers

import (
	"net/http"

	"github.com/garbein/lottery-golang/apps"
	"github.com/garbein/lottery-golang/routes"
)

type HttpServer struct {
	host string
	port string
}

func NewHttpServer() HttpServer {
	return HttpServer{host: apps.App.Config.ServerConfig.Host, port: apps.App.Config.ServerConfig.Port}
}

func (httpServer HttpServer) Run() {
	routes := routes.InitRoute()
	s := http.Server{
		Addr:    httpServer.host + ":" + httpServer.port,
		Handler: routes,
	}
	s.ListenAndServe()
}
