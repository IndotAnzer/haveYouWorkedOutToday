package router

import (
	"haveYouWorkedToday/controllers"
	"haveYouWorkedToday/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	auth := r.Group("/api/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)

	}

	api := r.Group("/api")
	api.Use(middlewares.AuthMiddleWare())
	{

		api.POST("/articles", controllers.CreateArticle)
		api.GET("/articles", controllers.GetAllArticles)
		api.GET("articles/:id", controllers.GetArticleByID)

		api.POST("/articles/:id/like", controllers.LikeArticle)
		api.GET("/articles/:id/like", controllers.GetArticleLikes)
	}
	return r
}
