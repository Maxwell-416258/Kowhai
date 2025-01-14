package cmd

import (
	"kowhai/apps/minio"
	"kowhai/configs"
	"kowhai/database"
	"kowhai/global"
	"kowhai/log"
	"kowhai/migrations"
)

func Start() {

	//初始化配置
	if global.Config == nil {
		global.Config = configs.InitConfig()
	}

	// 初始化日志组件
	global.Logger = log.InitLogger()

	// 初始化DB
	global.DB = database.InitDB()

	// 通过migrate来控制是否进行数据库migrate
	var migrate bool
	migrate = false

	if migrate {
		migrations.Migrate(global.DB)
	}

	//初始化minio相关数据
	_ = minio.InitMinioClient()
	_ = minio.InitStorageBuckets()
}
