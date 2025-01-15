package apps

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"kowhai/apps/user"
	"kowhai/apps/video"
	"time"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger()) //启用logger中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v1 := r.Group("/v1")
	{
		v1.POST("/user/create", user.CreateUser)
		v1.GET("/user/getbyname", user.GetUserByName)
		v1.GET("/users", user.GetUsers)
		v1.POST("/user/login", user.Login)
	}
	{
		v1.POST("/video/upload", video.UploadVedio)
	}
	return r
}
