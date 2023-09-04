package api

import (
	"github.com/gin-gonic/gin"
	"music-files/internal/config"
	"music-files/internal/handlers/cover_file_handlers"
	"music-files/internal/handlers/dir_handlers"
	"music-files/internal/handlers/music_file_handlers"
)

func SetupRouter(cfg *config.Configuration) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/music-files-service")
	{

		dirs := api.Group("/dirs")
		{
			dirs.POST("/add", dir_handlers.Add)
			dirs.DELETE("/remove", dir_handlers.Remove)
			dirs.POST("/scan", dir_handlers.Scan)
		}

		musicFiles := api.Group("/music-files")
		{
			musicFiles.GET("/{music_file_id}/download", music_file_handlers.GetFile)
			musicFiles.GET("/{music_file_id}", music_file_handlers.GetData)
			musicFiles.GET("/", music_file_handlers.GetDatas)
		}

		covers := api.Group("/covers")
		{
			covers.GET("/{cover_id}/download", cover_file_handlers.GetFile)
			covers.GET("/{cover_id}", cover_file_handlers.GetData)
			covers.GET("/", cover_file_handlers.GetDatas)
		}
	}

	return r
}
