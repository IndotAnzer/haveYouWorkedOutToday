package controllers

import (
	"haveYouWorkedOutToday/global"
	"haveYouWorkedOutToday/models"
	"haveYouWorkedOutToday/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateComment(ctx *gin.Context) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var comment models.Comment

	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment.UserID = userID

	articleID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	comment.ArticleID = uint(articleID)

	if err := global.Db.AutoMigrate(
		&models.Comment{},
		&models.Reply{},
	); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := global.Db.Create(&comment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := global.Db.Preload("User").First(&comment, comment.ID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"comment": comment})
}

func CreateReply(ctx *gin.Context) {
	var reply models.Reply

	if err := ctx.ShouldBindJSON(&reply); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := utils.GetUserID(ctx)

	reply.UserID = userID

	commentID, err := strconv.Atoi(ctx.Param("commentId"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	reply.CommentID = uint(commentID)

	if err := global.Db.AutoMigrate(&models.Reply{}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := global.Db.Create(&reply).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := global.Db.Preload("User").Preload("ParentReply").Preload("ParentReply.User").First(&reply, reply.ID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"reply": reply})
}

func GetCommentsAndReplies(ctx *gin.Context) {
	articleID := ctx.Param("id")

	var comments []models.Comment

	if err := global.Db.Preload("User").Preload("Replies").Preload("Replies.User").Preload("Replies.ParentReply").Preload("Replies.ParentReply.User").Where("article_id = ?", articleID).Find(&comments).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"comments": comments})
}

func DeleteComment(ctx *gin.Context) {
	commentID := ctx.Param("commentId")

	var comment models.Comment

	if err := global.Db.Where("id = ?", commentID).Preload("Replies").Unscoped().Delete(&models.Comment{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"comment": comment})
}

func DeleteReply(ctx *gin.Context) {
	replyID := ctx.Param("replyId")

	var reply models.Reply

	if err := global.Db.Where("id = ?", replyID).Unscoped().Delete(&models.Reply{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"reply": reply})
}
