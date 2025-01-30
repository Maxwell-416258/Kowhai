package video

import (
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"io"
	"kowhai/apps/minio"
	"kowhai/apps/user"
	"kowhai/bin"
	"kowhai/global"
	"math"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 上传视频(包括处理视频)
func UploadVedio(c *gin.Context) {
	// 限制文件大小，避免上传过大的文件
	const MaxUploadSize = 5000 << 20 // 1000MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadSize)

	Id := c.PostForm("userId")
	videoName := c.PostForm("videoName")
	if Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "userId 参数不能为空", "err": ""})
		return
	}

	// label
	labelStr := c.PostForm("label")
	if labelStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "label字段不能为空", "err": ""})
		return
	}

	// 将字符串转为整数
	labelInt, err := strconv.Atoi(labelStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "label必须为数字", "err": ""})
		return
	}

	// 转换为 VideoType 类型并验证
	label := VideoType(labelInt)
	if label < TypeMusic || label > TypeOther {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "label值不在有效范围内", "err": ""})
		return
	}

	userId, _ := strconv.Atoi(Id)

	// 获取上传的视频文件
	file, _, err := c.Request.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "视频文件获取失败",
			"err": err,
		})
		return
	}
	// 获取上传的视频封面
	image, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "视频封面获取失败",
			"err": err,
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
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "视频处理失败", "err": err})
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
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "视频封面保存失败", "err": err})
		return
	}

	// 保存视频封面链接到数据库
	imageLink := fmt.Sprintf("%s%s", minio_path, imageName)

	// 保存视频信息到数据库
	videoLink := fmt.Sprintf("%s%s", minio_path, m3u8)
	video := &Video{Name: videoName, UserId: userId, Link: videoLink, Image: imageLink, Label: label}
	if err = global.DB.Save(video).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "视频信息保存到数据库失败", "err": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "video save successful!", "data": nil})
	global.Logger.Info("video save successful!")
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
	if err := global.DB.Model(&Video{}).Where("id = ?", video_id).Select("sumLike").Error; err != nil {
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
	video_name := c.Query("name")
	if video_name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "name 参数不能为空", "err": ""})
		global.Logger.Error("name 参数不能为空")
		return
	}
	var video_list *[]Video
	if err := global.DB.Where("name like ?", "%"+video_name+"%").Order("create_time DESC").Find(&video_list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "视频搜索失败", "err": err})
		global.Logger.Error("视频搜索失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "搜索成功", "data": video_list})
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
		subscibeIds = append(subscibeIds, subscribe.Id)
	}

	var userInfo []map[string]interface{}
	//查询userName和Avatar
	err = global.DB.Model(&user.User{}).Select("user_name,avatar").Where("id IN ?", subscibeIds).Find(&userInfo).Error
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
