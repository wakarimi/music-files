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

// readResponse godoc
// @Description Response structure containing details of a single track.
// @Property TrackId (int) Unique ID of the track
// @Property CoverId (int, optional) Optional ID of the cover associated with the track
// @Property DurationMs (int64) Duration of the track in milliseconds
// @Property SizeByte (int64) Size of the track file in bytes
// @Property AudioCodec (string) Codec used for the audio track (e.g., "mp3", "flac")
// @Property BitrateKbps (int) Bitrate of the audio track in kilobits per second
// @Property SampleRateHz (int) Sample rate of the audio track in hertz
// @Property Channels (int) Number of audio channels (e.g., 2 for stereo)
// @Property HashSha256 (string) SHA-256 hash of the track file for integrity verification
type readResponse struct {
	TrackId      int    `json:"trackId"`
	CoverId      *int   `json:"coverId,omitempty"`
	DurationMs   int64  `json:"durationMs"`
	SizeByte     int64  `json:"sizeByte"`
	AudioCodec   string `json:"audioCodec"`
	BitrateKbps  int    `json:"bitrateKbps"`
	SampleRateHz int    `json:"sampleRateHz"`
	Channels     int    `json:"channels"`
	HashSha256   string `json:"hashSha256"`
}

// Read godoc
// @Summary Retrieve a single track by its ID
// @Tags Tracks
// @Accept json
// @Produce json
// @Param trackId path int true "Track ID"
// @Success 200 {object} readResponse "Successfully retrieved the track"
// @Failure 400 {object} types.ErrorResponse "Invalid trackId format"
// @Failure 500 {object} types.ErrorResponse "Failed to fetch track"
// @Router /tracks/{trackId} [get]
func (h *Handler) Read(c *gin.Context) {
	log.Debug().Msg("Fetching data about track")

	trackIdStr := c.Param("trackId")
	trackId, err := strconv.Atoi(trackIdStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid trackId format")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
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
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error: "Failed to fetch track",
		})
		return
	}

	log.Debug().Int("trackId", track.TrackId).Str("relativePath", track.RelativePath).Msg("Track fetched successfully")
	c.JSON(http.StatusOK, readResponse{
		TrackId:      track.TrackId,
		CoverId:      track.CoverId,
		DurationMs:   track.DurationMs,
		SizeByte:     track.SizeByte,
		AudioCodec:   track.AudioCodec,
		BitrateKbps:  track.BitrateKbps,
		SampleRateHz: track.SampleRateHz,
		Channels:     track.Channels,
		HashSha256:   track.HashSha256,
	})
}
