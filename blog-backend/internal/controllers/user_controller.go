package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"blog-backend/config"
	"blog-backend/internal/models"
)

type UpdateUserInput struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

// 获取当前用户信息
func GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "未授权访问", "error_code": "UNAUTHORIZED"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不存在", "error_code": "NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器内部错误", "error_code": "INTERNAL_ERROR"})
		return
	}

	// 隐藏密码字段
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// 更新用户信息
func UpdateCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "未授权访问", "error_code": "UNAUTHORIZED"})
		return
	}

	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入参数无效", "error_code": "INVALID_INPUT"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不存在", "error_code": "NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器内部错误", "error_code": "INTERNAL_ERROR"})
		return
	}

	// 更新用户信息
	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Avatar != "" {
		user.Avatar = input.Avatar
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新失败", "error_code": "INTERNAL_ERROR"})
		return
	}

	// 隐藏密码字段
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新成功",
		"data":    user,
	})
}
