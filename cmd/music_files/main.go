package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"music-files/api"
	"music-files/internal/config"
	"music-files/internal/context"
	"music-files/internal/database"
	"os"
)

// @title Wakarimi Music Files API
// @version 0.2

// @contact.name Dmitry Kolesnikov (Zalimannard)
// @contact.email zalimannard@mail.ru

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8023
// @BasePath /api/music-files-service
func main() {
	cfg := loadConfiguration()

	initializeLogger(cfg.Logger.Level)

	db := initializeDatabase(cfg.Database.ConnectionString)
	defer closeDatabase(db)
	initializeMigrations(db)

	ac := context.AppContext{
		Config: cfg,
		Db:     db,
	}

	server := initializeServer(&ac)
	runServer(server, ac.Config.HttpServer.Port)
}

func loadConfiguration() (cfg *config.Configuration) {
	cfg, err := config.LoadConfiguration()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to load configuration")
	}
	log.Debug().Msg("Configuration loaded")
	return cfg
}

func initializeLogger(level zerolog.Level) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).
		With().Caller().Logger().
		With().Str("service", "music_files").Logger().
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

func initializeServer(ac *context.AppContext) (r *gin.Engine) {
	r = api.SetupRouter(ac)
	log.Debug().Msg("Router initialized")
	return r
}

func runServer(server *gin.Engine, port string) {
	if err := server.Run(":" + port); err != nil {
		log.Panic().Err(err).Msg("Failed to start server")
	}
}
