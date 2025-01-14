package main

import (
	"kowhai/apps"
	"kowhai/cmd"
	"kowhai/global"
)

// @title 示例 API 文档
// @version 1.0
// @description 这是一个使用 Gin 框架的示例 Swagger 文档。
// @termsOfService http://example.com/terms/

// @contact.name API 支持团队
// @contact.url http://example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /v1
func main() {

	cmd.Start()

	// 路由
	r := apps.InitRouter()

	global.Logger.Info("Server started at :8081")

	r.Run(":8081")

}
