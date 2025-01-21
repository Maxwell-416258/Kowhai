package global

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"kowhai/config"
)

var Logger *zap.SugaredLogger

var Config *config.Config

var DB *gorm.DB

var Mongo *mongo.Client
