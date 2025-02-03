package main

import (
	"kowhai/apps/streaming"
	"kowhai/cmd"
	"kowhai/global"
	_ "net/http/pprof"
)

func main() {

	cmd.Start()

	// 路由
	r := streaming.InitRouter()

	global.Logger.Info("Server started at :8081")

	r.Run("0.0.0.0:8081")

}
