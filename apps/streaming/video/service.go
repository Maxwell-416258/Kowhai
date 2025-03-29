package video

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"io"
	minio2 "kowhai/apps/streaming/minio"
	"kowhai/apps/streaming/user"
	"kowhai/ffmpeg"
	"kowhai/global"
	"math"
	"net/http"
	"strconv"
)

// 上传视频(包括处理视频)
func UploadVideo(c *gin.Context) {
	// 限制文件大小
	const MaxUploadSize = 5000 << 20 // 1000MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadSize)

	Id := c.PostForm("userId")
	videoName := c.PostForm("videoName")
	if Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "userId 参数不能为空", "err": ""})
		return
	}

	// 处理 label
	labelStr := c.PostForm("label")
	if labelStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "label字段不能为空", "err": ""})
		return
	}
	labelInt, err := strconv.Atoi(labelStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "label必须为数字", "err": ""})
		return
	}
	label := VideoType(labelInt)
	if label < TypeMusic || label > TypeOther {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "label值不在有效范围内", "err": ""})
		return
	}

	userId, _ := strconv.Atoi(Id)

	// 获取上传的视频文件
	file, _, err := c.Request.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "视频文件获取失败", "err": err.Error()})
		return
	}
	defer file.Close()

	// 获取上传的视频封面
	image, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "视频封面获取失败", "err": err.Error()})
		return
	}
	defer image.Close()

	// 生成存储路径
	hlsDir := videoName
	m3u8 := fmt.Sprintf("%s.m3u8", videoName)
	ts := fmt.Sprintf("%s_%%03d.ts", videoName)
	minioPath := fmt.Sprintf("http://%s:%s/%s/%d/%s/", global.Config.Minio.Host, global.Config.Minio.Port, minio2.VEDIO_BUCKET, userId, hlsDir)

	// **立即返回响应，不等待视频处理**
	c.JSON(http.StatusOK, gin.H{"msg": "文件上传成功，正在处理中", "videoName": videoName})

	// **创建 Pipe**
	pr, pw := io.Pipe()

	// **异步处理**
	go func() {
		defer pw.Close()

		// 将文件内容写入 PipeWriter
		_, err := io.Copy(pw, file)
		if err != nil {
			global.Logger.Error("文件复制失败", err)
			return
		}
	}()

	// **异步启动 FFmpeg 处理**
	go func() {
		err := ffmpeg.Start(ts, m3u8, minioPath, hlsDir, userId, pr)
		if err != nil {
			global.Logger.Error("视频处理失败", err)
			return
		}

		// 上传封面到 MinIO
		imageName := fmt.Sprintf("%s.jpg", videoName)
		err = minio2.UploadVideo(userId, hlsDir, imageName, image)
		if err != nil {
			global.Logger.Error("视频封面保存失败", err)
			return
		}

		// 保存视频信息到数据库
		imageLink := fmt.Sprintf("%s%s", minioPath, imageName)
		videoLink := fmt.Sprintf("%s%s", minioPath, m3u8)
		video := &Video{Name: videoName, UserId: userId, Link: videoLink, Image: imageLink, Label: label}

		if err = global.DB.Save(video).Error; err != nil {
			global.Logger.Error("视频信息保存到数据库失败", err)
			return
		}

		global.Logger.Info("视频处理完成，已成功保存到数据库")
	}()
}

// 获取视频列表
func GetVideos(c *gin.Context) {
	var page, pageSize int
	var total int64
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ = strconv.Atoi(c.DefaultQuery("pageSize", "25"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 25
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
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "获取视频列表失败", "err": ""})
		return
	}

	data := map[string]interface{}{
		"videoList":  videoList,
		"page":       page,
		"pageSize":   pageSize,
		"total":      total,
		"totalPages": int(math.Ceil(float64(total) / float64(pageSize))),
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取成功",
		"data": data,
	})
}

// 获取用户上传的视频
func GetVideosByUserId(c *gin.Context) {
	var page, pageSize int
	var total int64
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ = strconv.Atoi(c.DefaultQuery("pageSize", "25"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 25
	}

	id := c.Query("userId")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "userId 参数不能为空", "err": ""})
		return
	}
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "userId必须为数字", "err": ""})
		return
	}

	var videoList []Video
	if err := global.DB.Where("user_id = ?", userId).Find(&videoList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "获取视频列表失败", "err": ""})
		return
	}

	data := map[string]interface{}{
		"videoList":  videoList,
		"page":       page,
		"pageSize":   pageSize,
		"total":      total,
		"totalPages": int(math.Ceil(float64(total) / float64(pageSize))),
	}

	c.JSON(http.StatusOK, gin.H{"msg": "获取成功", "data": data})
}

// 获取特定标签类型的视频
func GetVideosByLabel(c *gin.Context) {
	var page, pageSize int
	var total int64
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ = strconv.Atoi(c.DefaultQuery("pageSize", "25"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 25
	}

	label := c.Query("label")
	if label == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "label字段不能为空", "err": ""})
		return
	}
	label_int, err := strconv.Atoi(label)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "label必须为数字", "err": ""})
		return
	}

	var videoList []struct {
		Video
		UserName string `json:"user_name"`
		Avatar   string `json:"avatar"`
	}
	offset := (page - 1) * pageSize
	if err := global.DB.Model(&Video{}).
		Where("label = ?", label_int).
		Select("videos.*, users.user_name, users.avatar").
		Joins("left join users on videos.user_id = users.id").
		Count(&total).
		Limit(pageSize).Offset(offset).
		Find(&videoList).Error; err != nil {
		global.Logger.Error("Failed to get videos", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "获取视频列表失败", "err": ""})
		return
	}
	data := map[string]interface{}{
		"videoList":  videoList,
		"page":       page,
		"pageSize":   pageSize,
		"total":      total,
		"totalPages": int(math.Ceil(float64(total) / float64(pageSize))),
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取成功",
		"data": data,
	})
}

// 获取视频点赞数
func GetSumLikes(c *gin.Context) {
	var sum int64
	video_id := c.Query("video_id")
	if video_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "video_id 参数不能为空", "err": ""})
		return
	}
	id, _ := strconv.Atoi(video_id)
	if err := global.DB.Model(&Video{}).Where("id = ?", id).Select("sum_like").Find(&sum).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "点赞数查询失败", "err": ""})
	}
	c.JSON(http.StatusOK, gin.H{"msg": "获取成功", "data": sum})
}

// 添加点赞
func AddLikes(c *gin.Context) {
	video_id, err := strconv.Atoi(c.PostForm("video_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "添加失败", "err": err})
		return
	}
	if err := global.DB.Model(&Video{}).Where("id = ?", video_id).Update("sum_like", gorm.Expr("sum_like + ?", 1)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "添加失败", "err": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "点赞成功", "data": ""})
}

// 模糊搜索
func GetVideoByName(c *gin.Context) {
	type VideoWithUser struct {
		Video
		UserName string `json:"user_name"`
		Avatar   string `json:"avatar"`
	}

	video_name := c.Query("name")
	if video_name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "name 参数不能为空", "err": ""})
		global.Logger.Error("name 参数不能为空")
		return
	}

	var video_list []Video
	if err := global.DB.Where("name LIKE ?", "%"+video_name+"%").Order("create_time DESC").Find(&video_list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "视频搜索失败", "err": err})
		global.Logger.Error("视频搜索失败")
		return
	}

	var videoWithUserList []VideoWithUser

	// 遍历视频列表，查询每个视频对应的用户信息
	for _, video := range video_list {
		var user user.User
		if err := global.DB.Where("id = ?", video.UserId).First(&user).Error; err != nil {
			// 如果用户信息查询失败，可以设置为默认值
			global.Logger.Error("用户信息查询失败: ", err)
			continue
		}

		// 创建一个新的结构体，包含视频和用户信息
		videoWithUser := VideoWithUser{
			Video:    video,
			UserName: user.UserName,
			Avatar:   user.Avatar,
		}

		videoWithUserList = append(videoWithUserList, videoWithUser)
	}

	c.JSON(http.StatusOK, gin.H{"msg": "搜索成功", "data": videoWithUserList})
}

// 查询用户订阅的用户名和头像
func GetSubscribe(c *gin.Context) {
	user_id, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		global.Logger.Error("user_id error,err:", err)
		return
	}
	var subscribes []Subscribe
	err = global.DB.Model(Subscribe{}).Where("user_id = ?", user_id).Find(&subscribes).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败", "err": ""})
		global.Logger.Error("订阅数据查询失败", err)
		return
	}
	var subscibeIds []int
	for _, subscribe := range subscribes {
		subscibeIds = append(subscibeIds, subscribe.SubscribeId)
	}

	var userInfo []map[string]interface{}
	//查询userName和Avatar
	err = global.DB.Model(&user.User{}).Select("id,user_name,avatar").Where("id IN ?", subscibeIds).Find(&userInfo).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "用户名和头像获取失败", "err": ""})
		global.Logger.Error("用户名和头像获取失败", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "查询成功", "data": userInfo})
}

// 添加订阅
func CreateSubscribe(c *gin.Context) {
	var subscribe Subscribe
	err := c.BindJSON(&subscribe)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "绑定失败", "err": err})
		return
	}
	subscribe_id := subscribe.SubscribeId
	err = global.DB.Where("id = ?", subscribe_id).First(&user.User{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "订阅者id不存在", "err": ""})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "数据库查询失败", "err": err})
		return
	}

	if err = global.DB.Create(&subscribe).Error; err != nil {
		// 检查是否是唯一约束错误
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			// 错误 1062：Duplicate entry
			c.JSON(http.StatusBadRequest, gin.H{"msg": "已经订阅该用户，不能重复订阅", "err": ""})
			return
		}

		// 其他错误
		global.Logger.Error("订阅失败 ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "订阅失败", "err": ""})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "订阅成功", "data": subscribe})
}
