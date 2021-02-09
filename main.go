package main

import (
	"github.com/garbein/lottery-golang/apps"
	"github.com/garbein/lottery-golang/servers"
)

func main() {
	// 初始化app
	apps.InitApp()
	httpSever := servers.NewHttpServer()
	// 运行service
	httpSever.Run()
}
