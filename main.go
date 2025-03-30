package main

import (
	"github.com/gin-gonic/gin"
	"kowhai/apps/streaming"
	"kowhai/cmd"
	"kowhai/global"
	_ "net/http/pprof"
)

func main() {
	gin.ForceConsoleColor()

	cmd.Start()

	// 路由
	r := streaming.InitRouter()

	r.Static("/static", "./frontend/static")
	r.Static("/imgs", "./frontend/imgs")
	r.Static("/fonts", "./frontend/fonts")

	r.NoRoute(func(c *gin.Context) {
		c.File("frontend/index.html")
	})

	global.Logger.Info("Server started at :8081")

	r.Run("0.0.0.0:8081")

}
