package apps

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"vidspark/apps/user"
	_ "vidspark/docs"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	v1 := r.Group("/v1")
	{
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		v1.POST("/user", user.CreateUser)
	}
	return r
}
