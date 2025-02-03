package video

import (
	"kowhai/apps/streaming/base"
)

type VideoType int

const (
	TypeMusic     VideoType = iota + 1 // 音乐
	TypeDelicious                      // 美食
	TypeScape                          //风景
	TypeGame4                          //游戏
	TypeGhost                          //鬼畜
	TypeSports                         //运动
	TypeTravel                         //旅游
	TypeOther                          //其他
)

type Video struct {
	Id         int        `json:"id" gorm:"PrimaryKey;autoIncrement;comment:视频id"`
	UserId     int        `json:"userId" gorm:"index;comment:用户id"`
	Name       string     `json:"name" gorm:"comment:视频名称"`
	Image      string     `json:"image" gorm:"comment:视频封面"`
	Link       string     `json:"link" gorm:"comment:视频链接"`
	SumLike    int        `json:"sumLike" gorm:"default:0;comment:点赞数"`
	SumComment int        `json:"sumComment" gorm:"default:0;comment:评论数"`
	Label      VideoType  `json:":label" gorm:"type:int;default:8;comment:视频类型"`
	Audit      base.Audit `json:"audit" gorm:"embedded"`
}

type Subscribe struct {
	Id          int        `json:"id" gorm:"PrimaryKey;autoIncrement;comment:订阅id"`
	UserId      int        `json:"userId" gorm:"uniqueIndex:idx_user_subscribe;comment:用户id"`
	SubscribeId int        `json:"subscribeId" gorm:"uniqueIndex:idx_user_subscribe;comment:订阅者id"`
	Audit       base.Audit `json:"audit" gorm:"embedded"`
}

// 指定表名
func (Subscribe) TableName() string {
	return "subscribes"
}
