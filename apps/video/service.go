package video

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"kowhai/apps/minio"
	"kowhai/bin"
	"kowhai/global"
	"net/http"
	"strconv"
)

//上传视频(包括处理视频)

import (
	"sync"
)

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

	// 获取上传的文件
	file, _, err := c.Request.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("文件获取失败: %v", err),
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
	minio_path := fmt.Sprintf("http://%s:%s/%s/%v/%s", global.Config.Minio.Host, global.Config.Minio.Port, minio.VEDIO_BUCKET, userId, hlsDir)

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

	// 保存视频信息到数据库
	videoLink := fmt.Sprintf("%s/%s", minio_path, m3u8)
	video := &Video{Name: videoName, UserId: userId, Duration: videoDuration, Link: videoLink}
	if err = global.DB.Save(video).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("视频信息保存到数据库失败:%v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "video save successful!"})
	global.Logger.Info("video save successful!")
}
