package audio_file_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handler/responses"
	"music-files/internal/models"
	"net/http"
	"time"
)

// getAudioFilesResponseItem represents each audioFile item in the getAudioFiles response
type getAudioFilesResponseItem struct {
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

// getAudioFilesResponse is the response model for the GetAll API
type getAudioFilesResponse struct {
	// Array containing audioFile items
	AudioFiles []getAudioFilesResponseItem `json:"audioFiles"`
}

// GetAll retrieves all audioFiles
// @Summary Retrieve all audioFiles
// @Description Retrieves a list of all audioFiles in the system
// @Tags AudioFiles
// @Accept  json
// @Produce  json
// @Success 200 {object} getAudioFilesResponse
// @Failure 500 {object} responses.Error "Internal Server Error"
// @Router /audio-files [get]
func (h *Handler) GetAll(c *gin.Context) {
	log.Debug().Msg("Getting audioFiles")

	var audioFiles []models.AudioFile
	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		audioFiles, err = h.AudioFileService.GetAll(tx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get audioFiles")
		c.JSON(http.StatusInternalServerError, responses.Error{
			Message: "Failed to get audioFiles",
			Reason:  err.Error(),
		})
		return
	}

	audioFilesResponseItems := make([]getAudioFilesResponseItem, len(audioFiles))
	for i, audioFile := range audioFiles {
		audioFilesResponseItems[i] = getAudioFilesResponseItem{
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
		}
	}

	log.Debug().Msg("AudioFiles got successfully")
	c.JSON(http.StatusOK, getAudioFilesResponse{
		AudioFiles: audioFilesResponseItems,
	})
}
