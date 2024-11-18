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
// @Param user body user.User true "用户信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /user [post]
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
