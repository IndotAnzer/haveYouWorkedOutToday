package controllers

import (
	"errors"
	"fmt"
	"haveYouWorkedOutToday/global"
	"haveYouWorkedOutToday/models"
	"haveYouWorkedOutToday/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateArticle(ctx *gin.Context) {
	var article models.Article
	if err := ctx.ShouldBind(&article); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查必要字段
	if article.Title == "" || article.Content == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Title and content are required"})
		return
	}

	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	article.UserID = userID

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
		if global.RedisDB != nil {
			likes, err := global.RedisDB.Get(likeKey).Result()
			if err == nil {
				var likesInt int
				if _, err := fmt.Sscanf(likes, "%d", &likesInt); err == nil {
					articles[i].Likes = likesInt
				}
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
	if global.RedisDB != nil {
		likes, err := global.RedisDB.Get(likeKey).Result()
		if err == nil {
			var likesInt int
			if _, err := fmt.Sscanf(likes, "%d", &likesInt); err == nil {
				article.Likes = likesInt
			}
		}
	}

	ctx.JSON(http.StatusOK, article)
}

func GetArticleByUser(ctx *gin.Context) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var articles []models.Article

	if err := global.Db.Where("user_id = ?", userID).Find(&articles).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range articles {
		likeKey := fmt.Sprintf("article:%d:likes", articles[i].ID)
		if global.RedisDB != nil {
			likes, err := global.RedisDB.Get(likeKey).Result()
			if err == nil {
				var likesInt int
				if _, err := fmt.Sscanf(likes, "%d", &likesInt); err == nil {
					articles[i].Likes = likesInt
				}
			}
		}
	}

	ctx.JSON(http.StatusOK, articles)
}

func DeleteArticle(ctx *gin.Context) {
	id := ctx.Param("id")

	var article models.Article
	result := global.Db.Where("id = ?", id).Preload("FitnessActions.ActionGroups").Unscoped().Delete(&article)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Article deleted"})
}
