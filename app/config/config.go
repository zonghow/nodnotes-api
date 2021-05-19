package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func init() {
	loadCfg()
}

var C = struct {
	Gorm struct {
		Dsn         string `mapstructure:"dsn"`
		AutoMigrate bool   `mapstructure:"auto_migrate"`
	} `mapstructure:"gorm"`
	JWT struct {
		Signature string `mapstructure:"signature"`
		MaxAgeDay int    `mapstructure:"max_age_day"`
	} `mapstructure:"jwt"`
	Hashids struct {
		Salt      string `mapstructure:"salt"`
		Alphabet  string `mapstructure:"alphabet"`
		MinLength int    `mapstructure:"min_length"`
	} `mapstructure:"hashids"`
	Http struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"http"`
}{}

func loadCfg() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		} else {
			panic(fmt.Errorf("load configuration: %w", err))
		}
	}
	err = viper.Unmarshal(&C)
	if err != nil {
		panic(fmt.Errorf("load configuration: unmarshal: %w", err))
	}
}
