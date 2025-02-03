package streaming

import (
	"kowhai/apps/streaming/comment"
	"kowhai/apps/streaming/middleware"
	"kowhai/apps/streaming/user"
	"kowhai/apps/streaming/video"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger()) //启用logger中间件

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // 设置允许的源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 设置允许的请求方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // 设置允许的请求头
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},          // 设置可以暴露的响应头
		AllowCredentials: true,                                                // 是否允许凭据
		MaxAge:           12 * time.Hour,                                      // 缓存预检请求的结果
	}))
	r.POST("/user/login", user.Login)
	r.POST("/user/create", user.CreateUser)
	v1 := r.Group("/v1")
	v1.Use(middleware.JWTMiddleware())

	//user
	{

		v1.GET("/user/getbyname", user.GetUserByName)
		v1.GET("/users", user.GetUsers)
		v1.PATCH("/user/avatar", user.UploadAvatar)
	}

	//video
	{
		v1.POST("/video/upload", video.UploadVedio)
		v1.GET("/videos", video.GetVideos)
		v1.GET("/video/like", video.GetSumLikes)
		v1.PATCH("video/like", video.AddLikes)
		v1.GET("/video/search", video.GetVideoByName)
		v1.GET("/video/getVideosByLabel", video.GetVideosByLabel)
		v1.GET("/video/getVideosByUserId", video.GetVideosByUserId)
	}

	//comment
	{
		v1.POST("/comment/add", comment.AddComment)
		v1.GET("/comment/total", comment.GetCommentTotal)
		v1.GET("/comment/list", comment.GetCommentList)
	}
	//subscirbe
	{
		v1.GET("/subscribe", video.GetSubscribe)
		v1.POST("/subscribe", video.CreateSubscribe)
	}

	return r
}
