package vedio

import (
	"vidspark/apps/base"
)

type Video struct {
	Id         int        `json:"id" gorm:"PrimaryKey;autoIncrement;comment:视频id"`
	UserId     int        `json:"userId" gorm:"index;comment:用户id"`
	Link       string     `json:"link" gorm:"comment:视频链接"`
	SumLike    int        `json:"sumLike" gorm:"comment:点赞数"`
	SumComment int        `json:"sumComment" gorm:"comment:评论数"`
	Duration   int        `json:"duration" gorm:"comment:视频时长"`
	Audit      base.Audit `json:"audit" gorm:"embedded"`
}
