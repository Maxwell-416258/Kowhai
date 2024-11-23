package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
