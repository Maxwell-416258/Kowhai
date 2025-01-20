package apps

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"kowhai/apps/comment"
	"kowhai/apps/middleware"
	"kowhai/apps/user"
	"kowhai/apps/video"
	"time"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger()) //启用logger中间件

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://http://119.45.154.194:9001", "http://localhost:3000"}, // 设置允许的源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                    // 设置允许的请求方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},                    // 设置允许的请求头
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},                             // 设置可以暴露的响应头
		AllowCredentials: true,                                                                   // 是否允许凭据
		MaxAge:           12 * time.Hour,                                                         // 缓存预检请求的结果
	}))
	r.POST("/user/login", user.Login)

	v1 := r.Group("/v1")
	v1.Use(middleware.JWTMiddleware())

	{
		v1.POST("/user/create", user.CreateUser)
		v1.GET("/user/getbyname", user.GetUserByName)
		v1.GET("/users", user.GetUsers)
		v1.PATCH("/user/avatar", user.UploadAvatar)
	}

	{
		v1.POST("/video/upload", video.UploadVedio)
		v1.GET("/videos", video.GetVideos)
	}

	{
		v1.POST("/comment/add", comment.AddComment)
	}
	return r
}
