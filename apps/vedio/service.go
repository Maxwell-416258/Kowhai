package vedio

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"kowhai/bin"
	"kowhai/global"
	"net/http"
	"strconv"
)

//上传视频(包括处理视频)

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

	userId, _ := strconv.Atoi(Id)

	// 获取上传的文件
	file, _, err := c.Request.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("文件获取失败: %v", err),
		})
		return
	}
	defer file.Close()

	////获取视频时长
	//videoDuration, err := bin.GetVideoDuration(file)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"error": fmt.Sprintf("文件获取失败: %v", err),
	//	})
	//	return
	//}
	////重置文件指针
	//file.Seek(0, io.SeekCurrent)

	pr, pw := io.Pipe()
	// 开启协程，将上传的数据写入管道
	go func() {
		defer pw.Close()
		_, err := io.Copy(pw, file)
		if err != nil {
			global.Logger.Error("Failed to copy file to pipe", err)
		}
	}()

	// HLS 输出路径
	//hlsOutputDir := fmt.Sprintf("%d_%s", userId, videoName)
	hlsM3U8File := fmt.Sprintf("%s.m3u8", videoName)
	hlsSegmentPattern := fmt.Sprintf("%s_%%03d.ts", videoName)
	err = bin.Start(hlsSegmentPattern, hlsM3U8File, pr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("视频处理失败：%v", err)})
		return
	}

	// 上传 HLS 文件到 MinIO
	//hlsM3U8Key := fmt.Sprintf("%s/%s.m3u8", hlsOutputDir, videoName)
	//videoLink, err := minio.UploadVideo(hlsM3U8File, hlsM3U8Key)
	//c.JSON(http.StatusInternalServerError, gin.H{
	//	"error": fmt.Sprintf("上传 M3U8 到 MinIO 失败: %v", err),
	//})
	//return

	// 上传 TS 文件到 MinIO
	//for i := 0; ; i++ {
	//	tsFile := fmt.Sprintf("%s/%s_%03d.ts", hlsOutputDir, videoName, i)
	//	if _, err := os.Stat(tsFile); os.IsNotExist(err) {
	//		break
	//	}
	//	tsKey := fmt.Sprintf("%s/%s_%03d.ts", hlsOutputDir, videoName, i)
	//	if _, err := minio.UploadVideo(tsFile, tsKey); err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{
	//			"error": fmt.Sprintf("上传 TS 到 MinIO 失败: %v", err),
	//		})
	//		return
	//	}
	//}

	//保存视频信息到数据库

	video := Video{Name: videoName, UserId: userId}
	global.DB.Save(video)
}
