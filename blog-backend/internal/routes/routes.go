package routes

import (
	"github.com/gin-gonic/gin"

	"blog-backend/internal/controllers"
	"blog-backend/internal/middleware"
)

func SetupRoutes(r *gin.Engine) {
	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 认证相关接口
		auth := v1.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		// 用户相关接口
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/me", controllers.GetCurrentUser)
			users.PUT("/me", controllers.UpdateCurrentUser)
		}

		// 文章相关接口
		articles := v1.Group("/articles")
		{
			articles.GET("", controllers.GetArticles)
			articles.GET("/:id", controllers.GetArticle)
			
			// 需要认证的接口
			articles.Use(middleware.AuthMiddleware())
			{
				articles.POST("", controllers.CreateArticle)
				articles.PUT("/:id", controllers.UpdateArticle)
				articles.DELETE("/:id", controllers.DeleteArticle)
			}
		}

		// 评论相关接口
		comments := v1.Group("/articles/:id/comments")
		{
			comments.GET("", controllers.GetComments)
			comments.POST("", middleware.AuthMiddleware(), controllers.CreateComment)
		}
		
		// 删除评论接口
		v1.DELETE("/comments/:id", middleware.AuthMiddleware(), controllers.DeleteComment)
	}
}