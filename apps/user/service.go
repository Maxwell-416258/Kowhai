package user

import (
	"github.com/gin-gonic/gin"
	"kowhai/apps/middleware"
	"kowhai/apps/minio"
	"kowhai/global"
	"net/http"
)

func CreateUser(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Binding failed", "details": err.Error()})
		global.Logger.Error("Binding failed", err.Error())
		return
	}

	if err := global.DB.Where("name = ?", user.Name).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		global.Logger.Error("User already exists")
		return
	}

	if err := global.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		global.Logger.Error("Failed to create user", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
	global.Logger.Info("User created successfully", user)
}

func GetUserByName(c *gin.Context) {
	var user User
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		global.Logger.Error("Name is required")
		return
	}
	if err := global.DB.First(&user, "name = ?", name).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user", "details": err.Error()})
		global.Logger.Error("Failed to get user", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User retrived successfully", "user": user})
	global.Logger.Info("User retrived successfully", user)

}

func GetUsers(c *gin.Context) {
	var users []User
	if err := global.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users", "details": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": "User retrived successfully", "users": users})
}

func Login(c *gin.Context) {
	var user User
	var loginData struct {
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Binding failed", "details": err.Error()})
		global.Logger.Error("Binding failed", err.Error())
		return
	}
	if err := global.DB.Where("name = ?", loginData.Name).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users", "details": err.Error()})
		global.Logger.Error("Failed to get user", err.Error())
	}
	if user.Password != loginData.Password {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login", "details": "password is wrong"})
		global.Logger.Error("Failed to login", "password is wrong")
	}
	token, err := middleware.CreateToken(loginData.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login", "details": err.Error()})
		global.Logger.Error("Failed to login", err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"message": "Login successfully", "token": token, "user": user})
	global.Logger.Info("Login successfully", user)
}

func UploadAvatar(c *gin.Context) {
	var user User
	Id := c.PostForm("id")
	//查询用户name
	if err := global.DB.First(&user, "id = ?", Id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user", "details": err.Error()})
		global.Logger.Error("Failed to get user", err.Error())
		return
	}
	avatar, _, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "avatar must be required"})
		global.Logger.Error("avatar must be required")
		return
	}
	// 保存到minio
	avatar_name := user.Name + ".png"
	err = minio.Uploadavatar(user.Id, avatar_name, avatar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload avatar", "details": err.Error()})
		global.Logger.Error("Failed to upload avatar", err.Error())
		return
	}
	// 构造头像链接
	avatar_url := minio.GetAvatarUrl(user.Id, avatar_name)
	user.Avatar = avatar_url
	// 更新用户头像链接
	if err = global.DB.Model(&User{}).Where("id = ?", user.Id).Update("avatar", avatar_url).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user", "details": err.Error()})
		global.Logger.Error("Failed to update user", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Avatar uploaded successfully", "user": user})
}
