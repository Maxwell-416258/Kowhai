package main

import (
	"vidspark/apps"
	"vidspark/migrations"
	"vidspark/tools"
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
	// 通过migrate来控制是否进行数据库migrate
	var migrate bool
	migrate = false

	if migrate {
		db := tools.InitDB()
		migrations.Migrate(db)
	}

	// 路由
	r := apps.InitRouter()
	r.Run(":8081")

}
