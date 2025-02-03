package config

import (
	"github.com/BurntSushi/toml"
	"log"
	"sync"
)

var (
	once            sync.Once
	config_instance *Config
)

type Config struct {
	MySql    MySqlConfig    `toml:"mysql"`
	Mongo    MongoConfig    `toml:"mongo"`
	Minio    MinioConfig    `toml:"minio"`
	Log      LogConfig      `toml:"log"`
	Postgres PostgresConfig `toml:"postgres"`
}

type MySqlConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

type MongoConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

type MinioConfig struct {
	Host      string `toml:"host"`
	Port      string `toml:"port"`
	AccessKey string `toml:"access_key"`
	SecretKey string `toml:"secret_key"`
}

type LogConfig struct {
	MaxSize    int  `toml:"max_size"`
	MaxBackups int  `toml:"max_backups"`
	MaxAge     int  `toml:"max_age"`
	Compress   bool `toml:"compress"`
	Develop    bool `toml:"develop"`
}
type PostgresConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

func InitConfig() *Config {
	once.Do(func() {
		log.Println("Loading configuration...")
		var config Config
		if _, err := toml.DecodeFile("config/dev.toml", &config); err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}
		config_instance = &config
	})
	return config_instance
}
