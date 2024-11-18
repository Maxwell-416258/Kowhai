package user

import (
	"time"
	"vidspark/apps/base"
)

type User struct {
	Id        int        `json:"id" gorm:"PrimaryKey;autoIncrement"`
	Name      string     `json:"name" binding:"required"`
	Gender    string     `json:"gender" binding:"required"`
	Birth     string     `json:"birth" binding:"required"`
	Password  string     `json:"password" binding:"required"`
	Email     string     `json:"email" gorm:"default:null"`
	Phone     string     `json:"phone" binding:"required"`
	Followers int        `json:"followers" gorm:"column:num_followers;default:0"`
	Following int        `json:"following" gorm:"column:following_count;default:0"`
	LastLogin time.Time  `json:"last_login" gorm:"default:null"`
	Avator    string     `json:"avator" binding:"required"`
	State     int        `json:"state" gorm:"default:0"`
	Audit     base.Audit `json:"audit" gorm:"embedded"`
}
