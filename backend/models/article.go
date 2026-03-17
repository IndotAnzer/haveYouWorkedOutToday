package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title          string          `gorm:"not null" json:"title"`
	Content        string          `gorm:"not null" json:"content"`
	Preview        string          `gorm:"not null" json:"preview"`
	Likes          int             `gorm:"default:0" json:"likes"`
	FitnessActions []FitnessAction `gorm:"foreignKey:ArticleID;references:ID" json:"fitness_actions"`
}

type FitnessAction struct {
	gorm.Model
	ArticleID    uint                 `gorm:"not null" json:"article_id"`
	ActionName   string               `gorm:"not null" json:"action_name"`
	Remark       string               `json:"remark"`
	ActionGroups []FitnessActionGroup `gorm:"foreignKey:ActionID;references:ID" json:"action_groups"`
}

type FitnessActionGroup struct {
	gorm.Model
	ActionID   uint   `gorm:"not null" json:"action_id"`
	GroupIndex int    `gorm:"not null" json:"group_index"`
	Weight     int    `gorm:"not null" json:"weight"`
	RepNum     int    `gorm:"not null" json:"rep_num"`
	Remark     string `json:"remark"`
}
