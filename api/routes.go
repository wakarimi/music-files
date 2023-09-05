package api

import (
	"github.com/gin-gonic/gin"
	"music-files/internal/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/music-files-service")
	{

		dirs := api.Group("/dirs")
		{
			dirs.GET("/", handlers.DirGetAll)
			dirs.POST("/", handlers.DirAdd)
			dirs.DELETE("/:dirId", handlers.DirRemove)
			dirs.POST("/:dirId/scan", handlers.DirScan)
			dirs.POST("/scan-all", handlers.DirScanAll)
		}

		tracks := api.Group("/tracks")
		{
			tracks.GET("/:trackId", handlers.TrackGet)
			tracks.GET("/")
			tracks.GET("/:trackId/download")
		}

		covers := api.Group("/covers")
		{
			covers.GET("/:coverId", handlers.CoverGet)
			covers.GET("/:coverId/download", handlers.CoverDownload)
		}
	}

	return r
}
