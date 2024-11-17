package vedio

import (
	"vidspark/apps/base"
)

type Video struct {
	Id         int        `json:"id" gorm:"PrimaryKey;autoIncrement"`
	UserId     int        `json:"userId" gorm:"index"`
	Link       string     `json:"link"`
	SumLike    int        `json:"sumLike"`
	SumComment int        `json:"sumComment"`
	Duration   int        `json:"duration"`
	Audit      base.Audit `json:"audit" gorm:"embedded"`
}
