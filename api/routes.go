package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/config"
	"music-files/internal/database/repository"
	"music-files/internal/handlers/cover"
	"music-files/internal/handlers/directory"
	"music-files/internal/handlers/track"
	"music-files/internal/middleware"
)

func SetupRouter(httpServerConfig *config.HttpServer, db *sqlx.DB) *gin.Engine {
	log.Debug().Msg("Router setup")
	gin.SetMode(gin.ReleaseMode)

	coverRepo := repository.NewCoverRepository(db)
	dirRepo := repository.NewDirRepository(db)
	trackRepo := repository.NewTrackRepository(db)

	coverHandler := cover.NewHandler(coverRepo, dirRepo)
	dirHandler := directory.NewHandler(dirRepo, coverRepo, trackRepo)
	trackHandler := track.NewHandler(trackRepo)

	r := gin.New()
	r.Use(middleware.ZerologMiddleware(log.Logger))

	api := r.Group("/api/music-files-service")
	{
		dirs := api.Group("/dirs")
		{
			dirs.GET("/", dirHandler.ReadAll)
			dirs.POST("/", dirHandler.Create)
			dirs.DELETE("/:dirId", dirHandler.Delete)
			dirs.POST("/:dirId/scan", dirHandler.Scan)
			dirs.POST("/scan-all", dirHandler.ScanAll)
		}
		tracks := api.Group("/tracks")
		{
			tracks.GET("/:trackId", trackHandler.Read)
			tracks.GET("/", trackHandler.ReadAll)
			tracks.GET("/:trackId/download", trackHandler.Download)
		}
		covers := api.Group("/covers")
		{
			covers.GET("/:coverId", coverHandler.Read)
			covers.GET("/:coverId/download", coverHandler.Download)
		}
	}

	return r
}
