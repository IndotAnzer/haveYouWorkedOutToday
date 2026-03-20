package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	UserID    uint    `gorm:"not null" json:"user_id"`
	User      User    `gorm:"foreignkey:UserID" json:"user"`
	ArticleID uint    `gorm:"not null" json:"article_id"`
	Content   string  `gorm:"not null" json:"content"`
	Replies   []Reply `gorm:"foreignKey:CommentID" json:"replies"`
}

type Reply struct {
	gorm.Model
	UserID        uint   `gorm:"not null" json:"user_id"`
	User          User   `gorm:"foreignkey:UserID" json:"user"`
	CommentID     uint   `gorm:"not null" json:"comment_id"`
	Content       string `gorm:"not null" json:"content"`
	ParentReplyID uint   `json:"parent_reply_id"`
	ParentReply   *Reply `gorm:"foreignKey:ParentReplyID" json:"parent_reply"`
}
