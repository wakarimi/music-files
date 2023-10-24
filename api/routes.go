package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"music-files/internal/context"
	"music-files/internal/database/repository/cover_repo"
	"music-files/internal/database/repository/dir_repo"
	"music-files/internal/database/repository/song_repo"
	"music-files/internal/handler/dir_handler"
	"music-files/internal/middleware"
	"music-files/internal/service"
	"music-files/internal/service/cover_service"
	"music-files/internal/service/dir_service"
	"music-files/internal/service/song_service"
)

func SetupRouter(ac *context.AppContext) (r *gin.Engine) {
	log.Debug().Msg("Router setup")
	gin.SetMode(gin.ReleaseMode)

	r = gin.New()
	r.Use(middleware.ZerologMiddleware(log.Logger))

	coverRepo := cover_repo.NewRepository()
	songRepo := song_repo.NewRepository()
	dirRepo := dir_repo.NewRepository()
	txManager := service.NewTransactionManager(*ac.Db)

	coverService := cover_service.NewService(coverRepo)
	songService := song_service.NewService(songRepo)
	dirService := dir_service.NewService(dirRepo, *coverService, *songService)

	dirHandler := dir_handler.NewHandler(*dirService, txManager)

	api := r.Group("/api")
	{
		api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		roots := api.Group("/roots")
		{
			roots.GET("/", dirHandler.GetRoots)
			roots.POST("/", dirHandler.TrackRoot)
			roots.DELETE("/:dirId", dirHandler.UntrackRoot)
		}

		dirs := api.Group("/dirs")
		{
			dirs.POST("/:dirId/scan", dirHandler.Scan)
			dirs.POST("/scan-all", dirHandler.ScanAll)
		}

		songs := api.Group("/songs")
		{
			songs.GET("")
		}

		covers := api.Group("/covers")
		{
			covers.GET("")
		}
	}

	return r
}
