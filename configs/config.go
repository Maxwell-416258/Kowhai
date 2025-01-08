package configs

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
	MySql MySqlConfig `toml:"mysql"`
	Minio MinioConfig `toml:"minio"`
}

type MySqlConfig struct {
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

func InitConfig() *Config {
	once.Do(func() {
		log.Println("Loading configuration...")
		var config Config
		if _, err := toml.DecodeFile("configs/dev.toml", &config); err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}
		config_instance = &config
	})
	return config_instance
}
