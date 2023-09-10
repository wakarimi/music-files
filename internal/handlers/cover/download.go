package cover

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
	log.Debug().Msg("Downloading cover")

	coverIdStr := c.Param("coverId")
	coverId, err := strconv.Atoi(coverIdStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid coverId format")
		c.JSON(http.StatusBadRequest, types.Error{
			Error: "Invalid coverId format",
		})
		return
	}
	log.Debug().Int("coverId", coverId).Msg("Url parameter read successfully")

	cover, err := h.CoverRepo.Read(coverId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read cover")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to read cover",
		})
		return
	}
	log.Debug().Str("relativePath", cover.RelativePath).Msg("Cover read successfully")

	dir, err := h.DirRepo.Read(cover.DirId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read dir")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to read dir",
		})
		return
	}
	log.Debug().Str("path", dir.Path).Msg("Dir read successfully")

	absolutePath := filepath.Join(dir.Path, cover.RelativePath, cover.Filename)
	file, err := os.Open(absolutePath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open cover file")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to open cover file",
		})
		return
	}
	defer file.Close()
	log.Debug().Str("filename", file.Name()).Msg("File read successfully")

	log.Debug().Msg("Cover sent successfully")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+cover.Filename)
	c.File(absolutePath)
}
