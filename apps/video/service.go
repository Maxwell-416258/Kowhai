package video

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"kowhai/apps/minio"
	"kowhai/bin"
	"kowhai/global"
	"math"
	"net/http"
	"strconv"
)

import (
	"sync"
)

// 上传视频(包括处理视频)
func UploadVedio(c *gin.Context) {
	// 限制文件大小，避免上传过大的文件
	const MaxUploadSize = 5000 << 20 // 1000MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadSize)

	Id := c.PostForm("userId")
	videoName := c.PostForm("videoName")
	if Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId 参数不能为空"})
		return
	}

	videoDuration := c.PostForm("videoDuration")

	userId, _ := strconv.Atoi(Id)

	// 获取上传的视频文件
	file, _, err := c.Request.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("视频文件获取失败: %v", err),
		})
		return
	}
	// 获取上传的视频封面
	image, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("视频封面获取失败: %v", err),
		})
		return
	}

	// 使用 WaitGroup 来等待所有异步操作完成
	var wg sync.WaitGroup

	// 创建管道
	pr, pw := io.Pipe()

	// 开启协程，将上传的数据写入管道
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer pw.Close()

		_, err := io.Copy(pw, file)
		if err != nil {
			global.Logger.Error("Failed to copy file to pipe", err)
		}
	}()

	// 视频存放的minio文件夹
	hlsDir := fmt.Sprintf("%s", videoName)
	m3u8 := fmt.Sprintf("%s.m3u8", videoName)
	ts := fmt.Sprintf("%s_%%03d.ts", videoName)
	minio_path := fmt.Sprintf("http://%s:%s/%s/%v/%s/", global.Config.Minio.Host, global.Config.Minio.Port, minio.VEDIO_BUCKET, userId, hlsDir)

	// 使用另一个 goroutine 启动 FFmpeg 命令并上传文件
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := bin.Start(ts, m3u8, minio_path, hlsDir, userId, pr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("视频处理失败：%v", err)})
			return
		}
	}()

	// 等待所有异步操作完成后关闭文件
	wg.Wait()
	file.Close()

	// 保存视频封面到minio
	imageName := fmt.Sprintf("%s.jpg", videoName)
	err = minio.UploadVideo(userId, hlsDir, imageName, image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("视频封面保存失败：%v", err)})
		return
	}

	// 保存视频封面链接到数据库
	imageLink := fmt.Sprintf("%s%s", minio_path, imageName)

	// 保存视频信息到数据库
	videoLink := fmt.Sprintf("%s%s", minio_path, m3u8)
	video := &Video{Name: videoName, UserId: userId, Duration: videoDuration, Link: videoLink, Image: imageLink}
	if err = global.DB.Save(video).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("视频信息保存到数据库失败:%v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "video save successful!"})
	global.Logger.Info("video save successful!")
}

func GetVideos(c *gin.Context) {
	var page, pageSize int
	var total int64
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ = strconv.Atoi(c.DefaultQuery("pageSize", "15"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 15
	}

	var videoList []struct {
		Video
		UserName string `json:"user_name"`
		Avatar   string `json:"avatar"`
	}
	offset := (page - 1) * pageSize
	if err := global.DB.Model(&Video{}).
		Select("videos.*, users.user_name, users.avatar").
		Joins("left join users on videos.user_id = users.id").
		Count(&total).
		Limit(pageSize).Offset(offset).
		Find(&videoList).Error; err != nil {
		global.Logger.Error("Failed to get videos", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取视频列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":       videoList,
		"page":       page,
		"pageSize":   pageSize,
		"total":      total,
		"totalPages": int(math.Ceil(float64(total) / float64(pageSize))),
	})
}

func GetSumLikes(c *gin.Context) {
	var sum int64
	video_id := c.Query("video_id")
	if video_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "video_id 参数不能为空"})
		return
	}
	if err := global.DB.Model(&Video{}).Where("id = ?", video_id).Select("sumLike").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞数查询失败"})
	}
	c.JSON(http.StatusOK, gin.H{"sumLike": sum})
}

func AddLikes(c *gin.Context) {
	video_id, err := strconv.Atoi(c.PostForm("video_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	if err := global.DB.Model(&Video{}).Where("id = ?", video_id).Update("sum_like", gorm.Expr("sum_like + ?", 1)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "点赞成功"})
}

// 模糊搜索
func GetVideoByName(c *gin.Context) {
	video_name := c.Query("name")
	if video_name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name 参数不能为空"})
		global.Logger.Error("name 参数不能为空")
		return
	}
	var video_list *[]Video
	if err := global.DB.Where("name like ?", "%"+video_name+"%").Find(&video_list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "视频搜索失败"})
		global.Logger.Error("视频搜索失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": video_list})
}
