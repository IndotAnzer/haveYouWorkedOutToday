package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	UserID         uint            `gorm:"not null;index" json:"user_id"`
	Title          string          `gorm:"not null" json:"title"`
	Content        string          `gorm:"not null" json:"content"`
	Preview        string          `gorm:"not null" json:"preview"`
	Likes          int             `gorm:"default:0" json:"likes"`
	FitnessActions []FitnessAction `gorm:"constraint:OnDelete:CASCADE" json:"fitness_actions"`
}

type FitnessAction struct {
	gorm.Model
	ArticleID    uint                 `gorm:"not null" json:"article_id"`
	ActionName   string               `gorm:"not null" json:"action_name"`
	Remark       string               `json:"remark"`
	ActionGroups []FitnessActionGroup `gorm:"constraint:OnDelete:CASCADE" json:"action_groups"`
}

type FitnessActionGroup struct {
	gorm.Model
	FitnessActionID uint   `gorm:"not null" json:"fitness_action_id"`
	GroupIndex      int    `gorm:"not null" json:"group_index"`
	Weight          int    `gorm:"not null" json:"weight"`
	RepNum          int    `gorm:"not null" json:"rep_num"`
	Remark          string `json:"remark"`
}
