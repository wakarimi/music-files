package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"music-files/internal/api/http/controller"
	"music-files/internal/api/http/middleware"
	"music-files/internal/config"
)

type Router struct {
	router     *gin.Engine
	controller *controller.Controller
	log        *zerolog.Logger
}

func New(controller *controller.Controller, log *zerolog.Logger) *Router {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(middleware.ZerologMiddleware(*log))
	router.Use(middleware.ProduceLanguageMiddleware())
	return &Router{
		router:     router,
		controller: controller,
		log:        log,
	}
}

func (r *Router) RegisterRoutes() {
	api := r.router.Group("/api")
	{
		roots := api.Group("/roots")
		{
			roots.POST("", r.controller.AddRoot)
		}
	}
}

func (r *Router) StartHTTPServer(cfg config.HTTPConfig) error {
	err := r.router.Run(fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return errors.Wrap(err, "failed to start http server")
	}
	return nil
}
