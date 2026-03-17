package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title          string `gorm:"not null"`
	Content        string `gorm:"not null"`
	Preview        string `gorm:"not null"`
	Likes          int    `gorm:"default:0"`
	FitnessActions []FitnessAction
}

type FitnessAction struct {
	gorm.Model
	ArticleID    uint   `gorm:"not null"`
	ActionName   string `gorm:"not null"`
	Remark       string
	ActionGroups []FitnessActionGroup
}

type FitnessActionGroup struct {
	gorm.Model
	ActionID   uint `gorm:"not null"`
	GroupIndex int  `gorm:"not null"`
	Weight     int  `gorm:"not null"`
	RepNum     int  `gorm:"not null"`
	Remark     string
}
