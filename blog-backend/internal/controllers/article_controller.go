package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"blog-backend/config"
	"blog-backend/internal/models"
)

type CreateArticleInput struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdateArticleInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int `json:"total_pages"`
}

// 获取文章列表
func GetArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var articles []models.Article
	var total int64

	// 获取文章总数
	config.DB.Model(&models.Article{}).Count(&total)

	// 获取文章列表，预加载作者信息
	if err := config.DB.Preload("Author").Offset(offset).Limit(limit).Order("created_at DESC").Find(&articles).Error; err != nil {
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
			"articles":   articles,
			"pagination": Pagination{Page: page, Limit: limit, Total: total, TotalPages: totalPages},
		},
	})
}

// 获取文章详情
func GetArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的文章ID", "error_code": "INVALID_INPUT"})
		return
	}

	var article models.Article
	// 预加载作者信息
	if err := config.DB.Preload("Author").First(&article, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "文章不存在", "error_code": "NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器内部错误", "error_code": "INTERNAL_ERROR"})
		return
	}

	// 增加浏览量
	article.Views++
	config.DB.Save(&article)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    article,
	})
}

// 创建文章
func CreateArticle(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "未授权访问", "error_code": "UNAUTHORIZED"})
		return
	}

	var input CreateArticleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入参数无效", "error_code": "INVALID_INPUT"})
		return
	}

	article := models.Article{
		Title:    input.Title,
		Content:  input.Content,
		AuthorID: userID.(uint),
	}

	if err := config.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建失败", "error_code": "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "创建成功",
		"data":    article,
	})
}

// 更新文章
func UpdateArticle(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "未授权访问", "error_code": "UNAUTHORIZED"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的文章ID", "error_code": "INVALID_INPUT"})
		return
	}

	var article models.Article
	if err := config.DB.First(&article, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "文章不存在", "error_code": "NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器内部错误", "error_code": "INTERNAL_ERROR"})
		return
	}

	// 检查是否有权限更新文章
	if article.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "权限不足", "error_code": "FORBIDDEN"})
		return
	}

	var input UpdateArticleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入参数无效", "error_code": "INVALID_INPUT"})
		return
	}

	// 更新文章信息
	if input.Title != "" {
		article.Title = input.Title
	}
	if input.Content != "" {
		article.Content = input.Content
	}

	if err := config.DB.Save(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新失败", "error_code": "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新成功",
		"data":    article,
	})
}

// 删除文章
func DeleteArticle(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "未授权访问", "error_code": "UNAUTHORIZED"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的文章ID", "error_code": "INVALID_INPUT"})
		return
	}

	var article models.Article
	if err := config.DB.First(&article, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "文章不存在", "error_code": "NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器内部错误", "error_code": "INTERNAL_ERROR"})
		return
	}

	// 检查是否有权限删除文章
	if article.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "权限不足", "error_code": "FORBIDDEN"})
		return
	}

	if err := config.DB.Delete(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除失败", "error_code": "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
	})
}