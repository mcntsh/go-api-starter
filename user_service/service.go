package userserv

import (
	"fmt"
	"path"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config represents individual configuration definitions
// derrived from the `config.json` file.
type Config struct {
	HTTPAddress string `mapstructure:"http_address"`
	HTTPTimeout int    `mapstructure:"http_timeout"`

	DBDriver   string `mapstructure:"db_driver"`
	DBUser     string `mapstructure:"db_user"`
	DBPassword string `mapstructure:"db_password"`
	DBAddress  string `mapstructure:"db_address"`
	DBTable    string `mapstructure:"db_table"`
	dBString   string

	EncodingJWT string `mapstructure:"encoding_jwt"`
}

// Service represents an instance of this web service.
type Service struct {
	Config *Config
}

var serv = &Service{&Config{}}

// StartService creates a new Service instance and loads it
// with configuration from the config file.
func StartService() (*Service, error) {
	var err error

	// Load in the config
	_, curPath, _, _ := runtime.Caller(0)
	curDir := path.Dir(curPath)

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(curDir)

	err = viper.ReadInConfig()
	if err != nil {
		logrus.Fatal("Could not read config file")
		return nil, err
	}

	err = viper.Unmarshal(serv.Config)
	if err != nil {
		logrus.Fatal("Config could not be marshaled")
		return nil, err
	}

	serv.Config.dBString = fmt.Sprintf("%s:%s@tcp(%s)/%s",
		serv.Config.DBUser,
		serv.Config.DBPassword,
		serv.Config.DBAddress,
		serv.Config.DBTable,
	)

	return serv, nil
}
