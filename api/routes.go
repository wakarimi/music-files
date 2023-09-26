package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"music-files/internal/context"
	"music-files/internal/database/repository/cover_repo"
	"music-files/internal/database/repository/dir_repo"
	"music-files/internal/database/repository/track_repo"
	"music-files/internal/handler/dir_handler"
	"music-files/internal/middleware"
	"music-files/internal/service"
	"music-files/internal/service/cover_service"
	"music-files/internal/service/dir_service"
	"music-files/internal/service/track_service"
)

func SetupRouter(ac *context.AppContext) (r *gin.Engine) {
	log.Debug().Msg("Router setup")
	gin.SetMode(gin.ReleaseMode)

	r = gin.New()
	r.Use(middleware.ZerologMiddleware(log.Logger))

	coverRepo := cover_repo.NewRepository()
	trackRepo := track_repo.NewRepository()
	dirRepo := dir_repo.NewRepository()
	txManager := service.NewTransactionManager(*ac.Db)

	coverService := cover_service.NewService(coverRepo)
	trackService := track_service.NewService(trackRepo)
	dirService := dir_service.NewService(dirRepo, *coverService, *trackService)

	dirHandler := dir_handler.NewHandler(*dirService, txManager)

	api := r.Group("/api/music-files-service")
	{
		api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		dir := api.Group("dirs")
		{
			dir.POST("/", dirHandler.Create)
			dir.POST("/:dirId/scan", dirHandler.Scan)
		}
	}

	return r
}
