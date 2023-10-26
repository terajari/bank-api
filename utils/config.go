package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBSource             string        `mapstructure:"DB_SOURCE"`
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	HTTPServer           string        `mapstructure:"HTTP_SERVER"`
	TokenSymmtricKey     string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadConfig(filepath string) (config Config, err error) {
	viper.SetConfigFile(filepath)
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
