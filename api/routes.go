package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"music-files/internal/context"
	"music-files/internal/database/repository/dir_repo"
	"music-files/internal/handler/dir_handler"
	"music-files/internal/middleware"
	"music-files/internal/service"
	"music-files/internal/service/dir_service"
)

func SetupRouter(ac *context.AppContext) (r *gin.Engine) {
	log.Debug().Msg("Router setup")
	gin.SetMode(gin.ReleaseMode)

	r = gin.New()
	r.Use(middleware.ZerologMiddleware(log.Logger))

	dirRepo := dir_repo.NewRepository()
	txManager := service.NewTransactionManager(*ac.Db)

	dirService := dir_service.NewService(dirRepo)

	dirHandler := dir_handler.NewHandler(*dirService, txManager)

	api := r.Group("/api/music-files-service")
	{
		api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		dir := api.Group("dirs")
		{
			dir.POST("/", dirHandler.Create)
		}
	}

	return r
}
