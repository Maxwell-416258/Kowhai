package user

import (
	"vidspark/apps/base"
)

type User struct {
	Id       int        `json:"id" gorm:"PrimaryKey"`
	Name     string     `json:"name"`
	Password string     `json:"password"`
	audit    base.Audit `json:"audit"`
}
