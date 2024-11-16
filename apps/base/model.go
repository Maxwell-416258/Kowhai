package base

import "time"

type Audit struct {
	CreateTime time.Time `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"autoCreateTime"`
}
