package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}


type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	// 設定ファイルの検索パス
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.SetConfigName("config") // config.yaml, config.json など
	v.SetConfigType("yaml")

	// 環境変数を読み込む
	v.AutomaticEnv()
	
	// 環境変数のキーマッピングを設定
	v.BindEnv("server.port", "SERVER_PORT")
	v.BindEnv("database.host", "DATABASE_HOST")
	v.BindEnv("database.port", "DATABASE_PORT")
	v.BindEnv("database.user", "DATABASE_USER")
	v.BindEnv("database.password", "DATABASE_PASSWORD")
	v.BindEnv("database.dbname", "DATABASE_DBNAME")
	v.BindEnv("database.sslmode", "DATABASE_SSLMODE")

	// デフォルト値の設定
	v.SetDefault("server.port", "8080")
	v.SetDefault("database.path", "./feedapp.db")
	v.SetDefault("database.port", 5432)

	// 設定ファイルを読み込む
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using environment variables and defaults.")
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
