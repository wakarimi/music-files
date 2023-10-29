package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"music-files/internal/context"
	"music-files/internal/database/repository/audio_file_repo"
	"music-files/internal/database/repository/cover_repo"
	"music-files/internal/database/repository/dir_repo"
	"music-files/internal/handler/audio_file_handler"
	"music-files/internal/handler/cover_handler"
	"music-files/internal/handler/dir_handler"
	"music-files/internal/middleware"
	"music-files/internal/service"
	"music-files/internal/service/audio_file_service"
	"music-files/internal/service/cover_service"
	"music-files/internal/service/dir_service"
	"music-files/internal/service/file_processor_service"
)

func SetupRouter(ac *context.AppContext) (r *gin.Engine) {
	log.Debug().Msg("Router setup")
	gin.SetMode(gin.ReleaseMode)

	r = gin.New()
	r.Use(middleware.ZerologMiddleware(log.Logger))

	coverRepo := cover_repo.NewRepository()
	audioFileRepo := audio_file_repo.NewRepository()
	dirRepo := dir_repo.NewRepository()
	txManager := service.NewTransactionManager(*ac.Db)

	coverService := cover_service.NewService(coverRepo)
	audioFileService := audio_file_service.NewService(audioFileRepo)
	dirService := dir_service.NewService(dirRepo, *coverService, *audioFileService)
	fileProcessorService := file_processor_service.NewService(*dirService, *coverService, *audioFileService)

	coverHandler := cover_handler.NewHandler(*coverService, *fileProcessorService, txManager)
	audioFileHandler := audio_file_handler.NewHandler(*audioFileService, *fileProcessorService, txManager)
	dirHandler := dir_handler.NewHandler(*dirService, txManager)

	api := r.Group("/api")
	{
		api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		roots := api.Group("/roots")
		{
			roots.GET("/", dirHandler.GetRoots)
			roots.POST("/", dirHandler.AddRootToWatchList)
			roots.DELETE("/:dirId", dirHandler.RemoveRootFromWatchList)
		}

		dirs := api.Group("/dirs")
		{
			dirs.GET("/:dirId", dirHandler.GetDir)
			dirs.GET("/:dirId/content", dirHandler.Content)
			dirs.POST("/:dirId/scan", dirHandler.Scan)
			dirs.POST("/scan", dirHandler.ScanAll)
		}

		audioFiles := api.Group("/audio-files")
		{
			audioFiles.GET("/:audioFileId", audioFileHandler.GetAudioFile)
			audioFiles.GET("/", audioFileHandler.GetAll)
			audioFiles.GET("/:audioFileId/download", audioFileHandler.Download)
			audioFiles.GET("/:audioFileId/cover", audioFileHandler.GetCover)
			audioFiles.GET("/sha256/:sha256", audioFileHandler.SearchBySha256)
		}

		covers := api.Group("/covers")
		{
			covers.GET("/:coverId", coverHandler.GetCover)
			covers.GET("/:coverId/download", coverHandler.Download)
		}
	}

	log.Debug().Msg("Router setup successfully")
	return r
}
