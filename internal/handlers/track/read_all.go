package track

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handlers/types"
	"music-files/internal/models"
	"net/http"
)

// readAllResponseItem godoc
// @Description Single track item structure in the response of ReadAll endpoint.
// @Property TrackId (int) Unique ID of the track.
// @Property CoverId (int, optional) Optional ID of the cover associated with the track.
// @Property DurationMs (int64) Duration of the track in milliseconds.
// @Property AudioCodec (string) Codec used for the audio track (e.g., "mp3", "flac").
// @Property Size (int64) Size of the track file in bytes.
// @Property HashSha256 (string) SHA-256 hash of the track file for integrity verification.
type readAllResponseItem struct {
	TrackId    int    `json:"trackId"`
	CoverId    *int   `json:"coverId,omitempty"`
	DurationMs int64  `json:"durationMs"`
	AudioCodec string `json:"audioCodec"`
	Size       int64  `json:"size"`
	HashSha256 string `json:"hashSha256"`
}

// readAllResponse godoc
// @Description Response structure containing details of all tracks.
// @Property Tracks (array of readAllResponseItem) Array containing details of all tracks.
type readAllResponse struct {
	Tracks []readAllResponseItem `json:"tracks"`
}

// ReadAll godoc
// @Summary Retrieve all tracks
// @Tags Tracks
// @Accept json
// @Produce json
// @Success 200 {object} readAllResponse "Successfully retrieved all tracks"
// @Failure 500 {object} types.ErrorResponse "Failed to fetch all tracks"
// @Router /tracks [get]
func (h *Handler) ReadAll(c *gin.Context) {
	log.Debug().Msg("Fetching all tracks")

	var tracks []models.Track

	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		tracks, err = h.TrackService.ReadAll(tx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch all tracks")
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error: "Failed to fetch all tracks",
		})
		return
	}

	tracksResponse := make([]readAllResponseItem, 0)
	for _, track := range tracks {
		trackResponse := readAllResponseItem{
			TrackId:    track.TrackId,
			CoverId:    track.CoverId,
			DurationMs: track.DurationMs,
			AudioCodec: track.AudioCodec,
			Size:       track.Size,
			HashSha256: track.HashSha256,
		}
		tracksResponse = append(tracksResponse, trackResponse)
	}

	log.Info().Int("tracksCount", len(tracks)).Msg("Track fetched successfully")
	c.JSON(http.StatusOK, readAllResponse{
		Tracks: tracksResponse,
	})
}
