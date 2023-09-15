package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"music-files/internal/context"
	"music-files/internal/database/repository"
	"music-files/internal/handlers/cover"
	"music-files/internal/handlers/directory"
	"music-files/internal/handlers/track"
	"music-files/internal/middleware"
	"music-files/internal/service"
	"music-files/internal/service/cover_service"
	"music-files/internal/service/dir_service"
	"music-files/internal/service/track_service"
)

func SetupRouter(ac *context.AppContext) (r *gin.Engine) {
	log.Debug().Msg("Router setup")
	gin.SetMode(gin.ReleaseMode)

	txManager := service.NewTransactionManager(*ac.Db)

	coverRepo := repository.NewCoverRepository(ac.Db)
	dirRepo := repository.NewDirRepository(ac.Db)
	trackRepo := repository.NewTrackRepository(ac.Db)

	coverService := cover_service.NewService(coverRepo, dirRepo)
	dirService := dir_service.NewService(dirRepo, coverRepo, trackRepo)
	trackService := track_service.NewService(trackRepo, dirRepo)

	coverHandler := cover.NewHandler(txManager, *coverService, *dirService)
	dirHandler := directory.NewHandler(txManager, *dirService, *coverService, *trackService)
	trackHandler := track.NewHandler(txManager, *trackService, *dirService)

	r = gin.New()
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
