package video

import (
	"kowhai/apps/base"
)

type Video struct {
	Id         int        `json:"id" gorm:"PrimaryKey;autoIncrement;comment:视频id"`
	UserId     int        `json:"userId" gorm:"index;comment:用户id"`
	Name       string     `json:"name" gorm:"comment:视频名称"`
	Link       string     `json:"link" gorm:"comment:视频链接"`
	SumLike    int        `json:"sumLike" gorm:"default:0;comment:点赞数"`
	SumComment int        `json:"sumComment" gorm:"default:0;comment:评论数"`
	Duration   string     `json:"duration" gorm:"comment:视频时长"`
	Audit      base.Audit `json:"audit" gorm:"embedded"`
}
