package config

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"time"
)

type AppConfig struct {
	LoggingLevel zerolog.Level
}

type HTTPConfig struct {
	Port int
}

type DBConfig struct {
	Host          string
	Port          int
	DBName        string
	User          string
	Password      string
	Timeout       time.Duration
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	Charset       string
	MigrationPath string
}

type Config struct {
	App  AppConfig
	HTTP HTTPConfig
	DB   DBConfig
}

func New() (config Config, err error) {
	viper.SetDefault("APP_LOGGING_LEVEL", "INFO")
	viper.SetDefault("HTTP_PORT", 8021)
	viper.SetDefault("DB_READ_TIMEOUT", "1s")
	viper.SetDefault("DB_WRITE_TIMEOUT", "1s")
	viper.SetDefault("DB_CHARSET", "UTF-8")

	viper.AutomaticEnv()

	config = Config{
		App: AppConfig{
			LoggingLevel: parseLoggingLevel(viper.GetString("APP_LOGGING_LEVEL")),
		},

		HTTP: HTTPConfig{
			Port: viper.GetInt("HTTP_PORT"),
		},

		DB: DBConfig{
			Host:          viper.GetString("DB_HOST"),
			Port:          viper.GetInt("DB_PORT"),
			DBName:        viper.GetString("DB_NAME"),
			User:          viper.GetString("DB_USER"),
			Password:      viper.GetString("DB_PASSWORD"),
			ReadTimeout:   viper.GetDuration("DB_READ_TIMEOUT"),
			WriteTimeout:  viper.GetDuration("DB_WRITE_TIMEOUT"),
			Charset:       viper.GetString("DB_CHARSET"),
			MigrationPath: "internal/storage/migration",
		},
	}

	return config, nil
}

func parseLoggingLevel(loggingLevel string) zerolog.Level {
	switch loggingLevel {
	case "TRACE":
		return zerolog.TraceLevel
	case "DEBUG":
		return zerolog.DebugLevel
	case "INFO":
		return zerolog.InfoLevel
	case "WARN":
		return zerolog.WarnLevel
	case "ERROR":
		return zerolog.ErrorLevel
	case "FATAL":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}
