package user

import (
	"time"
	"vidspark/apps/base"
	"vidspark/apps/vedio"
)

type User struct {
	Id        int           `json:"id" gorm:"PrimaryKey;autoIncrement"`
	Name      string        `json:"name"`
	Password  string        `json:"password"`
	Email     string        `json:"email" gorm:"default:null"`
	Phone     string        `json:"phone"`
	LastLogin time.Time     `json:"last_login"`
	Avator    string        `json:"avator"`
	Vedios    []vedio.Video `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	State     int           `json:"state" gorm:"default:0"`
	Audit     base.Audit    `json:"audit"`
}
