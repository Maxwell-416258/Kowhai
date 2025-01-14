package global

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"kowhai/configs"
)

var Logger *zap.SugaredLogger

var Config *configs.Config

var DB *gorm.DB
