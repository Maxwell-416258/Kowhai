package comment

import (
	"kowhai/apps/base"
	"time"
)

type Comment struct {
	Id          int        `json:"id" gorm:"autoIncrement;comment:评论id"`
	VideoId     string     `json:"video_id" gorm:"comment:视频id"`
	ReviewerId  string     `json:"reviewer_id" gorm:"comment:评论者id"`
	Content     string     `json:"content" gorm:"comment:评论内容"`
	CommentTime time.Time  `json:"comment_time" gorm:"comment:评论时间"`
	Likes       int        `json:"likes" gorm:"comment:评论点赞数"`
	Audit       base.Audit `json:"audit" gorm:"embedded"`
}
