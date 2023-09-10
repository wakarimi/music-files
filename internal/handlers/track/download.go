package track

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"music-files/internal/handlers/types"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func (h *Handler) Download(c *gin.Context) {
	log.Debug().Msg("Downloading track")

	trackIdStr := c.Param("trackId")
	trackId, err := strconv.Atoi(trackIdStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid trackId format")
		c.JSON(http.StatusBadRequest, types.Error{
			Error: "Invalid trackId format",
		})
		return
	}
	log.Debug().Int("trackId", trackId).Msg("Url parameter read successfully")

	track, err := h.TrackRepo.Read(trackId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read track")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to read track",
		})
		return
	}
	log.Debug().Str("relativePath", track.RelativePath).Msg("Track read successfully")

	dir, err := h.DirRepo.Read(track.DirId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read dir")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to read dir",
		})
		return
	}
	log.Debug().Str("path", dir.Path).Msg("Dir read successfully")

	absolutePath := filepath.Join(dir.Path, track.RelativePath, track.Filename)
	file, err := os.Open(absolutePath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open track file")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to open track file",
		})
		return
	}
	defer file.Close()
	log.Debug().Str("filename", file.Name()).Msg("File read successfully")

	log.Debug().Msg("Track sent successfully")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+track.Filename)
	c.File(absolutePath)
}
