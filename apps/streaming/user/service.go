package user

import (
	"github.com/gin-gonic/gin"
	"kowhai/apps/streaming/middleware"
	"kowhai/apps/streaming/minio"
	"kowhai/global"
	"net/http"
)

func CreateUser(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Binding failed", "err": err.Error()})
		global.Logger.Error("Binding failed", err.Error())
		return
	}

	if err := global.DB.Where("name = ?", user.UserName).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "User already exists", "err": ""})
		global.Logger.Error("User already exists")
		return
	}

	if err := global.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to create user", "err": ""})
		global.Logger.Error("Failed to create user", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "User created successfully",
		"data": user,
	})
	global.Logger.Info("User created successfully", user)
}

func GetUserByName(c *gin.Context) {
	var user User
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Name is required", "err": ""})
		global.Logger.Error("Name is required")
		return
	}
	if err := global.DB.First(&user, "name = ?", name).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to get user", "err": err.Error()})
		global.Logger.Error("Failed to get user", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "User retrived successfully", "data": user})
	global.Logger.Info("User retrived successfully", user)

}

func GetUsers(c *gin.Context) {
	var users []User
	if err := global.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to get users", "err": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"msg": "User retrived successfully", "data": users})
}

func Login(c *gin.Context) {
	var user User
	var loginData struct {
		Name     string `json:"user_name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.BindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Binding failed", "err": err.Error()})
		global.Logger.Error("Binding failed", err.Error())
		return
	}
	if err := global.DB.Where("user_name = ?", loginData.Name).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to get users", "err": err.Error()})
		global.Logger.Error("Failed to get user", err.Error())
		return
	}
	if user.Password != loginData.Password {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to login", "err": "password is wrong"})
		global.Logger.Error("Failed to login", "password is wrong")
		return
	}
	token, err := middleware.CreateToken(loginData.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to login", "err": err.Error()})
		global.Logger.Error("Failed to login", err.Error())
		return
	}

	data := map[string]interface{}{
		"token": token,
		"user":  user,
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Login successfully", "data": data})
	global.Logger.Info("Login successfully", user)
}

func UploadAvatar(c *gin.Context) {
	var user User
	Id := c.PostForm("id")
	//查询用户name
	if err := global.DB.First(&user, "id = ?", Id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to get user", "err": err.Error()})
		global.Logger.Error("Failed to get user", err.Error())
		return
	}
	avatar, _, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "avatar must be required", "err": ""})
		global.Logger.Error("avatar must be required")
		return
	}
	// 保存到minio
	avatar_name := user.UserName + ".png"
	err = minio.Uploadavatar(user.Id, avatar_name, avatar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to upload avatar", "err": err.Error()})
		global.Logger.Error("Failed to upload avatar", err.Error())
		return
	}
	// 构造头像链接
	avatar_url := minio.GetAvatarUrl(user.Id, avatar_name)
	user.Avatar = avatar_url
	// 更新用户头像链接
	if err = global.DB.Model(&User{}).Where("id = ?", user.Id).Update("avatar", avatar_url).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to update user", "err": err.Error()})
		global.Logger.Error("Failed to update user", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Avatar uploaded successfully", "data": user})
}
