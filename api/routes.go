package api

import (
	"github.com/gin-gonic/gin"
	"music-files/internal/config"
	"music-files/internal/handlers"
)

func SetupRouter(cfg *config.Configuration) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/music-files-service")
	{
		api.POST("/scan", func(c *gin.Context) { handlers.Scan(c, cfg) })
	}

	return r
}
