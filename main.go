package main

import (
	"kowhai/apps"
	"kowhai/cmd"
	"kowhai/global"
	_ "net/http/pprof"
)

func main() {

	cmd.Start()

	// 路由
	r := apps.InitRouter()

	global.Logger.Info("Server started at :8081")

	r.Run(":8081")

}
