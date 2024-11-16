package comment

import (
	"vidspark/apps/base"
)

type Comment struct {
	Id         int        `json:"id" gorm:"autoIncrement"`
	VideoId    string     `json:"video_id"`
	ReviewerId string     `json:"reviewer_id"`
	Content    string     `json:"content"`
	Audit      base.Audit `json:"audit"`
}
