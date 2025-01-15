package user

import (
	"github.com/gin-gonic/gin"
	"kowhai/global"
	"net/http"
)

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
	if user.Password == loginData.Password {
		c.JSON(http.StatusOK, gin.H{"message": "Login successfully", "user": user})
		global.Logger.Info("Login successfully", user)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login", "details": "password is wrong"})
		global.Logger.Error("Failed to login", "password is wrong")
	}
}
