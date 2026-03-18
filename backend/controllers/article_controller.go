package controllers

import (
	"errors"
	"fmt"
	"haveYouWorkedOutToday/global"
	"haveYouWorkedOutToday/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateArticle(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user models.User
	if err := global.Db.Where("username = ?", username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	var article models.Article
	if err := ctx.ShouldBind(&article); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article.UserID = user.ID

	if err := global.Db.AutoMigrate(
		&models.Article{},
		&models.FitnessAction{},
		&models.FitnessActionGroup{},
	); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := global.Db.Create(&article).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, article)
}

func GetAllArticles(ctx *gin.Context) {
	var articles []models.Article

	if err := global.Db.Find(&articles).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range articles {
		likeKey := fmt.Sprintf("article:%d:likes", articles[i].ID)
		likes, err := global.RedisDB.Get(likeKey).Result()
		if err == nil {
			var likesInt int
			if _, err := fmt.Sscanf(likes, "%d", &likesInt); err == nil {
				articles[i].Likes = likesInt
			}
		}
	}

	ctx.JSON(http.StatusOK, articles)
}

func GetArticleByID(ctx *gin.Context) {
	id := ctx.Param("id")

	var article models.Article

	err := global.Db.Preload("FitnessActions.ActionGroups").Where("id = ?", id).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	likeKey := fmt.Sprintf("article:%d:likes", article.ID)
	likes, err := global.RedisDB.Get(likeKey).Result()
	if err == nil {
		var likesInt int
		if _, err := fmt.Sscanf(likes, "%d", &likesInt); err == nil {
			article.Likes = likesInt
		}
	}

	ctx.JSON(http.StatusOK, article)
}

func GetArticleByUser(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user models.User

	if err := global.Db.Where("username = ?", username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var articles []models.Article

	if err := global.Db.Where("user_id = ?", user.ID).Find(&articles).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range articles {
		likeKey := fmt.Sprintf("article:%d:likes", articles[i].ID)
		likes, err := global.RedisDB.Get(likeKey).Result()
		if err == nil {
			var likesInt int
			if _, err := fmt.Sscanf(likes, "%d", &likesInt); err == nil {
				articles[i].Likes = likesInt
			}
		}
	}

	ctx.JSON(http.StatusOK, articles)
}

func DeleteArticle(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := global.Db.Where("id = ?", id).Preload("FitnessActions.ActionGroups").Unscoped().Delete(&models.Article{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Article deleted"})
}
