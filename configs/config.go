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
}

type MySqlConfig struct {
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
		if _, err := toml.DecodeFile("configs/dev.toml", &config); err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}
		config_instance = &config
	})
	return config_instance
}
