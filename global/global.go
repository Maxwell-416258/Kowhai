package global

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"vidspark/configs"
)

var Logger *zap.SugaredLogger

var Config *configs.Config

var DB *gorm.DB
