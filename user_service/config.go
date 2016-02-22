package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	HTTPAddress string `mapstructure:"http_address"`
	HTTPTimeout int    `mapstructure:"http_timeout"`

	DBDriver   string `mapstructure:"db_driver"`
	DBUser     string `mapstructure:"db_user"`
	DBPassword string `mapstructure:"db_password"`
	DBAddress  string `mapstructure:"db_address"`
	DBTable    string `mapstructure:"db_table"`
	DBString   string

	EncodingJWT string `mapstructure:"encoding_jwt"`
}

var config *Config

func configLoad() error {
	config = &Config{}

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(config)
	if err != nil {
		return err
	}

	config.DBString = fmt.Sprintf("%s:%s@tcp(%s)/%s", config.DBUser, config.DBPassword, config.DBAddress, config.DBTable)

	return nil
}
