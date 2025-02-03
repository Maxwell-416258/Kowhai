package cmd

import (
	"kowhai/apps/streaming/minio"
	"kowhai/config"
	"kowhai/database"
	"kowhai/global"
	"kowhai/log"
	"kowhai/migration"
)

func Start() {

	//初始化配置
	if global.Config == nil {
		global.Config = config.InitConfig()
	}

	// 初始化日志组件
	global.Logger = log.InitLogger()

	// 初始化DB
	global.DB = database.InitDB()

	// 通过migrate来控制是否进行数据库migrate

	// 初始化mongo
	global.Mongo = database.InitMongo()

	// 初始化postgres
	var err error
	global.DbPool, err = database.InitPostgres()
	if err != nil {
		global.Logger.Fatalf("Failed to connect to postgres: %v", err)
	}

	if MIGRATE {
		migration.Migrate(global.DB)
	}

	//初始化minio相关数据
	_ = minio.InitMinioClient()
	_ = minio.InitStorageBuckets()
}
