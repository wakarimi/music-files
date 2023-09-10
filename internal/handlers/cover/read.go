package cover

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"music-files/internal/handlers/types"
	"net/http"
	"strconv"
)

type readResponse struct {
	CoverId   int    `json:"coverId"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
}

func (h *Handler) Read(c *gin.Context) {
	log.Debug().Msg("Fetching data about cover")

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

	log.Debug().Int("coverId", cover.CoverId).Str("relativePath", cover.RelativePath).Msg("Cover fetched successfully")
	c.JSON(http.StatusOK, readResponse{
		CoverId:   cover.CoverId,
		Extension: cover.Extension,
		Size:      cover.Size,
	})
}
