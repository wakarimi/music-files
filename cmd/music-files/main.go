package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"music-files/api"
	"music-files/internal/config"
	"music-files/internal/database"
	"os"
)

func main() {
	cfg := loadConfiguration()
	initializeLogger(cfg.Logger.Level)

	db := initializeDatabase(cfg.Database.ConnectionString)
	defer closeDatabase(db)
	initializeMigrations(db)

	server := initializeServer(&cfg.HttpServer, db)
	runServer(server, cfg.HttpServer.Port)
}

func loadConfiguration() *config.Configuration {
	cfg, err := config.LoadConfiguration()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to load configuration")
	}
	log.Debug().Msg("Configuration loaded")
	return cfg
}

func initializeLogger(level zerolog.Level) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).
		//With().Caller().Logger().
		With().Str("service", "music-files").Logger().
		Level(level)
	log.Debug().Msg("Logger initialized")
}

func initializeDatabase(connectionString string) (db *sqlx.DB) {
	db, err := database.ConnectDb(connectionString)
	if err != nil {
		log.Panic().Msg("Failed to initialize database")
	}
	log.Debug().Msg("Database initialized")
	return db
}

func closeDatabase(db *sqlx.DB) {
	if err := db.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close database connection")
	}
	log.Debug().Msg("Database connection closed")
}

func initializeMigrations(db *sqlx.DB) {
	if err := database.RunMigrations(db, "./internal/database/migrations"); err != nil {
		log.Panic().Err(err).Msg("Failed to apply migrations")
	}
	log.Debug().Msg("Data schema actualized")
}

func initializeServer(httpServerConfig *config.HttpServer, db *sqlx.DB) *gin.Engine {
	r := api.SetupRouter(httpServerConfig, db)
	log.Debug().Msg("Router initialized")
	return r
}

func runServer(server *gin.Engine, port string) {
	if err := server.Run(":" + port); err != nil {
		log.Panic().Err(err).Msg("Failed to start server")
	}
}
