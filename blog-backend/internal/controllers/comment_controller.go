package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"blog-backend/config"
	"blog-backend/internal/models"
)

type CreateCommentInput struct {
	Content string `json:"content" binding:"required"`
}

// 获取文章评论列表
func GetComments(c *gin.Context) {
	articleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的文章ID", "error_code": "INVALID_INPUT"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var comments []models.Comment
	var total int64

	// 获取评论总数
	config.DB.Model(&models.Comment{}).Where("article_id = ?", articleID).Count(&total)

	// 获取评论列表，预加载作者信息
	if err := config.DB.Preload("Author").Where("article_id = ?", articleID).Offset(offset).Limit(limit).Order("created_at DESC").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器内部错误", "error_code": "INTERNAL_ERROR"})
		return
	}

	totalPages := int(total)/limit + 1
	if int(total)%limit == 0 {
		totalPages = int(total) / limit
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"comments":   comments,
			"pagination": Pagination{Page: page, Limit: limit, Total: total, TotalPages: totalPages},
		},
	})
}

// 发表评论
func CreateComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "未授权访问", "error_code": "UNAUTHORIZED"})
		return
	}

	articleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的文章ID", "error_code": "INVALID_INPUT"})
		return
	}

	// 检查文章是否存在
	var article models.Article
	if err := config.DB.First(&article, articleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "文章不存在", "error_code": "NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器内部错误", "error_code": "INTERNAL_ERROR"})
		return
	}

	var input CreateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入参数无效", "error_code": "INVALID_INPUT"})
		return
	}

	comment := models.Comment{
		Content:   input.Content,
		ArticleID: uint(articleID),
		AuthorID:  userID.(uint),
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "评论发表失败", "error_code": "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "评论发表成功",
		"data":    comment,
	})
}

// 删除评论
func DeleteComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "未授权访问", "error_code": "UNAUTHORIZED"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的评论ID", "error_code": "INVALID_INPUT"})
		return
	}

	var comment models.Comment
	if err := config.DB.First(&comment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "评论不存在", "error_code": "NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器内部错误", "error_code": "INTERNAL_ERROR"})
		return
	}

	// 检查是否有权限删除评论 (评论作者或文章作者)
	var article models.Article
	config.DB.First(&article, comment.ArticleID)

	if comment.AuthorID != userID.(uint) && article.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "权限不足", "error_code": "FORBIDDEN"})
		return
	}

	if err := config.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除失败", "error_code": "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
	})
}