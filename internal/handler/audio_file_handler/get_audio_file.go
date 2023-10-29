package audio_file_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/handler/responses"
	"music-files/internal/models"
	"net/http"
	"strconv"
	"time"
)

// getAudioFileResponse represents audioFile in the getAudioFile response
type getAudioFileResponse struct {
	// Unique identifier for the audioFile
	AudioFileId int `json:"audioFileId"`
	// Directory identifier where the audioFile resides
	DirId int `json:"dirId"`
	// Filename of the audioFile
	Filename string `json:"filename"`
	// File extension of the audioFile
	Extension string `json:"extension"`
	// File size in bytes
	SizeByte int64 `json:"sizeByte"`
	// Duration of the audioFile in milliseconds
	DurationMs int64 `json:"durationMs"`
	// Bitrate in kilobits per second
	BitrateKbps int `json:"bitrateKbps"`
	// Sample rate in hertz
	SampleRateHz int `json:"sampleRateHz"`
	// Number of audio channels
	ChannelsN int `json:"channelsN"`
	// SHA-256 hash of the file
	Sha256 string `json:"sha256"`
	// Time of the last update to the audioFile's content
	LastContentUpdate time.Time `json:"lastContentUpdate"`
}

// GetAudioFile retrieves a audioFile by its identifier
// @Summary Retrieve a audioFile by its ID
// @Description Retrieves a single audioFile by its ID
// @Tags AudioFiles
// @Accept  json
// @Produce  json
// @Param   audioFileId path int true "AudioFile ID"
// @Success 200 {object} getAudioFileResponse
// @Failure 400,404,500 {object} responses.Error
// @Router /audio-files/{audioFileId} [get]
func (h *Handler) GetAudioFile(c *gin.Context) {
	log.Debug().Msg("Getting audioFile")

	audioFileIdStr := c.Param("audioFileId")
	audioFileId, err := strconv.Atoi(audioFileIdStr)
	if err != nil {
		log.Error().Err(err).Str("audioFileIdStr", audioFileIdStr).Msg("Invalid audioFileId format")
		c.JSON(http.StatusBadRequest, responses.Error{
			Message: "Invalid audioFileId format",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("audioFileId", audioFileId).Msg("Url parameter read successfully")

	var audioFile models.AudioFile
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		audioFile, err = h.AudioFileService.GetAudioFile(tx, audioFileId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get audioFile")
		if _, ok := err.(errors.NotFound); ok {
			c.JSON(http.StatusNotFound, responses.Error{
				Message: "AudioFile not found",
				Reason:  err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, responses.Error{
				Message: "Failed to get audioFile",
				Reason:  err.Error(),
			})
		}
		return
	}

	log.Debug().Msg("AudioFile got successfully")
	c.JSON(http.StatusOK, getAudioFileResponse{
		AudioFileId:       audioFile.AudioFileId,
		DirId:             audioFile.DirId,
		Filename:          audioFile.Filename,
		Extension:         audioFile.Extension,
		SizeByte:          audioFile.SizeByte,
		DurationMs:        audioFile.DurationMs,
		BitrateKbps:       audioFile.BitrateKbps,
		SampleRateHz:      audioFile.SampleRateHz,
		ChannelsN:         audioFile.ChannelsN,
		Sha256:            audioFile.Sha256,
		LastContentUpdate: audioFile.LastContentUpdate,
	})
}
