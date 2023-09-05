package config

import (
	"github.com/spf13/viper"
	"strings"
)

type Configuration struct {
	DatabaseConfiguration
	HttpServerConfiguration
}

type DatabaseConfiguration struct {
	DatabaseConnectionString string
}

type HttpServerConfiguration struct {
	Port string
}

func LoadConfiguration() (config *Configuration, err error) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	config = &Configuration{
		DatabaseConfiguration{
			DatabaseConnectionString: viper.GetString("WAKARIMI_MUSIC_FILES_DB_STRING"),
		},
		HttpServerConfiguration{
			Port: viper.GetString("HTTP_SERVER_PORT"),
		},
	}

	return config, nil
}
