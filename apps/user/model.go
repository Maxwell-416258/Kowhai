package user

import (
	"kowhai/apps/base"
	"time"
)

type User struct {
	Id        int        `json:"id" gorm:"PrimaryKey;autoIncrement;comment:用户id"`
	UserName  string     `json:"user_name" binding:"required" gorm:"comment:用户名"`
	Gender    string     `json:"gender" binding:"required" gorm:"comment:性别"`
	Birth     string     `json:"birth" binding:"required" gorm:"comment:出生日期"`
	Password  string     `json:"password" binding:"required" gorm:"comment:密码"`
	Email     string     `json:"email" gorm:"default:null;comment:邮箱"`
	Phone     string     `json:"phone" binding:"required" gorm:"comment:电话号码"`
	Followers int        `json:"followers" gorm:"column:num_followers;default:0;comment:关注你的"`
	Following int        `json:"following" gorm:"column:following_count;default:0;comment:你关注的"`
	LastLogin time.Time  `json:"last_login" gorm:"default:null;comment:最后一次登录的时间"`
	Avatar    string     `json:"avatar" gorm:"default:null;comment:头像链接"`
	State     int        `json:"state" gorm:"default:0;comment:状态"`
	Audit     base.Audit `json:"audit" gorm:"embedded"`
}
