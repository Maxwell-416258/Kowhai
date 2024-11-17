package user

import (
	"time"
	"vidspark/apps/base"
)

type User struct {
	Id        int        `json:"id" gorm:"PrimaryKey;autoIncrement"`
	Name      string     `json:"name"`
	Gender    string     `json:"gender"`
	Birth     time.Time  `json:"birth"`
	Password  string     `json:"password"`
	Email     string     `json:"email" gorm:"default:null"`
	Phone     string     `json:"phone"`
	Followers int        `json:"followers" gorm:"column:num_followers"`
	Following int        `json:"following" gorm:"column:following_count"`
	LastLogin time.Time  `json:"last_login"`
	Avator    string     `json:"avator"`
	State     int        `json:"state" gorm:"default:0"`
	Audit     base.Audit `json:"audit" gorm:"embedded"`
}
