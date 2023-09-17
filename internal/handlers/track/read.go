package track

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handlers/types"
	"music-files/internal/models"
	"net/http"
	"strconv"
)

type readResponse struct {
	TrackId    int    `json:"trackId"`
	CoverId    *int   `json:"coverId,omitempty"`
	AudioCodec string `json:"audioCodec"`
	Size       int64  `json:"size"`
	HashSha256 string `json:"hashSha256"`
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

	var track models.Track

	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		track, err = h.TrackService.Read(tx, trackId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch track")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to fetch track",
		})
		return
	}

	log.Debug().Int("trackId", track.TrackId).Str("relativePath", track.RelativePath).Msg("Track fetched successfully")
	c.JSON(http.StatusOK, readResponse{
		TrackId:    track.TrackId,
		CoverId:    track.CoverId,
		AudioCodec: track.AudioCodec,
		Size:       track.Size,
		HashSha256: track.HashSha256,
	})
}
