package controllers

import (
	"errors"
	"haveYouWorkedOutToday/global"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func LikeArticle(ctx *gin.Context) {
	articleID := ctx.Param("id")
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	likeKey := "article:" + articleID + ":likes"
	userLikeKey := "article:" + articleID + ":user:" + username.(string)

	liked, err := global.RedisDB.Get(userLikeKey).Result()
	if err == nil && liked == "1" {
		if err := global.RedisDB.Decr(likeKey).Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := global.RedisDB.Del(userLikeKey).Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Successfully unliked article", "action": "unliked"})
		return
	}

	if err := global.RedisDB.Incr(likeKey).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := global.RedisDB.Set(userLikeKey, "1", 0).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully liked article", "action": "liked"})
}

func GetArticleLikes(ctx *gin.Context) {
	articleID := ctx.Param("id")
	username, exists := ctx.Get("username")

	likeKey := "article:" + articleID + ":likes"

	likes, err := global.RedisDB.Get(likeKey).Result()
	if errors.Is(err, redis.Nil) {
		likes = "0"
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	isLiked := false
	if exists {
		userLikeKey := "article:" + articleID + ":user:" + username.(string)
		liked, err := global.RedisDB.Get(userLikeKey).Result()
		if err == nil && liked == "1" {
			isLiked = true
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"count": likes, "is_liked": isLiked})
}
