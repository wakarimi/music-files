package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"music-files/internal/context"
	"music-files/internal/middleware"
)

func SetupRouter(ac *context.AppContext) (r *gin.Engine) {
	log.Debug().Msg("Router setup")
	gin.SetMode(gin.ReleaseMode)

	r = gin.New()
	r.Use(middleware.ZerologMiddleware(log.Logger))

	//coverRepo := cover_repo.NewRepository()
	//trackRepo := track_repo.NewRepository()
	//dirRepo := dir_repo.NewRepository()
	//txManager := service.NewTransactionManager(*ac.Db)

	//coverService := cover_service.NewService(coverRepo)
	//trackService := track_service.NewService(trackRepo)
	//dirService := dir_service.NewService(dirRepo, *coverService, *trackService)

	//dirHandler := dir_handler.NewHandler(*dirService, txManager)

	api := r.Group("/api")
	{
		api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		dirs := api.Group("dirs")
		{
			dirs.GET("")
		}

		tracks := api.Group("tracks")
		{
			tracks.GET("")
		}

		covers := api.Group("covers")
		{
			covers.GET("")
		}
	}

	return r
}
