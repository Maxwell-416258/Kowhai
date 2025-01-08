package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"vidspark/apps/minio"
	"vidspark/tools"
)

var db = tools.InitDB()

// CreateUser 创建用户
// @Summary 创建用户
// @Description 创建一个新的用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param name body string true "Name"
// @Param gender body string true "Gender"
// @Param birth body string true "Birth"
// @Param password body string true "Password"
// @Param email body string false "Email"
// @Param phone body string true "Phone"
// @Param avator body string true "Avator"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /user/create [post]
func CreateUser(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Binding failed", "details": err.Error()})
		return
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

// GetUserByName 查询用户
// @Summary 查询用户
// @Description 查询用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param name query string true "用户名"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /user/getbyname [get]
func GetUserByName(c *gin.Context) {
	var user User
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}
	if err := db.First(&user, "name = ?", name).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User retrived successfully", "user": user})

}

// GetUsers 查询所有用户
// @Summary 查询所有用户
// @Description 查询所有用户
// @Tags 用户
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /users [get]
func GetUsers(c *gin.Context) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users", "details": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": "User retrived successfully", "users": users})
}

// 用户上传视频接口（使用 Gin）
func UploadVideoHandler(c *gin.Context) {
	// 限制文件大小，避免上传过大的文件
	const MaxUploadSize = 50 << 20 // 10MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadSize)

	Id := c.PostForm("userId")

	if Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId 参数不能为空"})
		return
	}
	userId, _ := strconv.Atoi(Id)

	// 获取上传的文件
	file, _, err := c.Request.FormFile("vedio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("文件获取失败: %v", err),
		})
		return
	}
	defer file.Close()

	// 获取文件名，可以根据需求生成唯一的文件名
	fileName := c.DefaultPostForm("fileName", "default_video.mp4")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "文件名不能为空",
		})
		return
	}

	// 调用 UploadVideo 方法上传文件
	videoURL, err := minio.UploadVideo(userId, file, fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("视频上传失败: %v", err),
		})
		return
	}

	// 返回成功上传的视频 URL
	c.JSON(http.StatusOK, gin.H{
		"message": "视频上传成功",
		"url":     videoURL,
	})
}
