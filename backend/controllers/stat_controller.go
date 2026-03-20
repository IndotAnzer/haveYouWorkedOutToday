package controllers

import (
	"haveYouWorkedOutToday/global"
	"haveYouWorkedOutToday/models"
	"haveYouWorkedOutToday/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTrainingFrequency(ctx *gin.Context) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	var stats []struct {
		Date  string `json:"date"`
		Count string `json:"count"`
	}

	query := global.Db.Model(&models.Article{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("user_id = ?", userID)

	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate+" 23:59:59")
	} else if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	} else if endDate != "" {
		query = query.Where("created_at <= ?", endDate+" 23:59:59")
	}

	if err := query.Group("DATE(created_at)").Find(&stats).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

func GetTrainingVolume(ctx *gin.Context) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	type Action struct {
		ActionName string  `json:"action_name"`
		Num        int     `json:"num"`
		Load       float64 `json:"load"`
	}

	var stats []struct {
		Date   string   `json:"date"`
		Action []Action `json:"action"`
	}

	var results []struct {
		Date        string  `json:"date"`
		ActionName  string  `json:"action_name"`
		TotalGroups int     `json:"total_groups"`
		TotalLoad   float64 `json:"total_load"`
	}

	query := global.Db.Model(&models.Article{}).
		Joins("LEFT JOIN fitness_actions ON fitness_actions.article_id = articles.id").
		Joins("LEFT JOIN fitness_action_groups ON fitness_action_groups.fitness_action_id = fitness_actions.id").
		Where("articles.user_id = ?", userID)

	if startDate != "" {
		query = query.Where("articles.created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("articles.created_at <= ?", endDate+" 23:59:59")
	}

	if err := query.Select("DATE_FORMAT(articles.created_at, '%Y-%m-%d') as date, fitness_actions.action_name as action_name, COUNT(fitness_action_groups.id) as total_groups, SUM(fitness_action_groups.weight) as total_load").
		Group("DATE_FORMAT(articles.created_at, '%Y-%m-%d'), fitness_actions.action_name").
		Find(&results).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dataMap := make(map[string][]Action)

	for _, result := range results {
		if result.ActionName != "" {
			dataMap[result.Date] = append(dataMap[result.Date], Action{
				ActionName: result.ActionName,
				Num:        result.TotalGroups,
				Load:       result.TotalLoad,
			})
		}
	}

	for date, actions := range dataMap {
		stats = append(stats, struct {
			Date   string   `json:"date"`
			Action []Action `json:"action"`
		}{
			Date:   date,
			Action: actions,
		})
	}

	ctx.JSON(http.StatusOK, stats)
}
