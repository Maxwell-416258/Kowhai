package base

import "time"

type Audit struct {
	CreateTime time.Time `json:"create_time" gorm:"autoCreateTime;comment:创建时间"`
	UpdateTime time.Time `json:"update_time" gorm:"autoCreateTime;comment:更新时间"`
}
