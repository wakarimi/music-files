package track

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"music-files/internal/handlers/types"
	"net/http"
	"strconv"
)

type readResponse struct {
	TrackId   int    `json:"trackId"`
	CoverId   int    `json:"coverId,omitempty"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
}

func (h *Handler) Read(c *gin.Context) {
	log.Debug().Msg("Fetching data about track")

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

	log.Debug().Int("trackId", track.TrackId).Str("relativePath", track.RelativePath).Msg("Track fetched successfully")
	c.JSON(http.StatusOK, readResponse{
		TrackId:   track.TrackId,
		CoverId:   *track.CoverId,
		Extension: track.Extension,
		Size:      track.Size,
	})
}
