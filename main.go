package main

import (
	"github.com/garbein/lottery-golang/apps"
	"github.com/garbein/lottery-golang/servers"
)

func main() {
	apps.InitApp()
	httpSever := servers.NewHttpServer()
	httpSever.Run()
}
