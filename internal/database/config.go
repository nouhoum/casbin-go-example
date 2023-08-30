package database

import (
	"github.com/samber/do"
	"github.com/spf13/viper"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
	SSLMode  string
	TimeZone string
}

func NewConfig(i *do.Injector) (*Config, error) {
	return &Config{
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
		User:     viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),
		Dbname:   viper.GetString("DB_NAME"),
		SSLMode:  viper.GetString("DB_SSL_MODE"),
		TimeZone: viper.GetString("DB_TIME_ZONE"),
	}, nil
}

func (cfg Config) DSN() string {
	return "host=" + cfg.Host +
		" user=" + cfg.User +
		" password=" + cfg.Password +
		" dbname=" + cfg.Dbname +
		" port=" + cfg.Port +
		" sslmode=" + cfg.SSLMode +
		" TimeZone=" + cfg.TimeZone
}
