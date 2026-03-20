package router

import (
	"haveYouWorkedOutToday/controllers"
	"haveYouWorkedOutToday/middlewares"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 开发环境允许所有来源，生产环境请指定具体域名
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
		api.GET("/my/articles", controllers.GetArticleByUser)
		api.GET("/articles/:id", controllers.GetArticleByID)
		api.DELETE("/articles/:id", controllers.DeleteArticle)

		api.POST("/articles/:id/like", controllers.LikeArticle)
		api.GET("/articles/:id/like", controllers.GetArticleLikes)

		api.POST("/articles/:id/comments", controllers.CreateComment)
		api.GET("/articles/:id/comments", controllers.GetCommentsAndReplies)
		api.POST("/articles/:id/comments/:commentId/replies", controllers.CreateReply)
		api.DELETE("/articles/:id/comments/:commentId", controllers.DeleteComment)
		api.DELETE("/articles/:id/comments/:commentId/replies/:replyId", controllers.DeleteReply)

		api.GET("/articles/stats/frequency", controllers.GetTrainingFrequency)
		api.GET("/articles/stats/volume", controllers.GetTrainingVolume)
	}
	return r
}
