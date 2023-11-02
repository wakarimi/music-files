package audio_file_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handler/response"
	"music-files/internal/model"
	"net/http"
	"time"
)

// searchBySha256ResponseItem represents a single audioFile item in the search by SHA256 response.
type searchBySha256ResponseItem struct {
	// Unique identifier for the audioFile.
	AudioFileId int `json:"audioFileId"`
	// Directory ID where the audioFile is located.
	DirId int `json:"dirId"`
	// Filename of the audioFile.
	Filename string `json:"filename"`
	// File extension of the audioFile.
	Extension string `json:"extension"`
	// File size of the audioFile in bytes.
	SizeByte int64 `json:"sizeByte"`
	// Duration of the audioFile in milliseconds.
	DurationMs int64 `json:"durationMs"`
	// Bitrate of the audioFile in Kbps.
	BitrateKbps int `json:"bitrateKbps"`
	// Sample rate of the audioFile in Hz.
	SampleRateHz int `json:"sampleRateHz"`
	// Number of channels in the audioFile.
	ChannelsN int `json:"channelsN"`
	// SHA-256 hash of the audioFile.
	Sha256 string `json:"sha256"`
	// Timestamp of the last content update.
	LastContentUpdate time.Time `json:"lastContentUpdate"`
}

// searchBySha256Response represents the search by SHA256 API response.
type searchBySha256Response struct {
	// Array of audioFiles that match the search query.
	AudioFiles []searchBySha256ResponseItem `json:"audioFiles"`
}

// SearchBySha256 retrieves a list of audioFiles based on SHA256 hash.
// @Summary Search audioFiles by SHA256 hash
// @Description Retrieves a list of audioFiles that have the specified SHA256 hash.
// @Tags AudioFiles
// @Accept  json
// @Produce  json
// @Param   sha256     path    string  true        "SHA256 hash"
// @Success 200 {object} searchBySha256Response
// @Failure 500 {object} response.Error "Internal Server Error"
// @Router /audio-files/sha256/{sha256} [get]
func (h *Handler) SearchBySha256(c *gin.Context) {
	sha256 := c.Param("sha256")
	log.Debug().Str("sha256", sha256).Msg("Url parameter read successfully")

	var audioFiles []model.AudioFile
	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		audioFiles, err = h.AudioFileService.SearchBySha256(tx, sha256)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error{
			Message: "Failed to get audioFiles",
			Reason:  err.Error(),
		})
		return
	}

	audioFilesResponseItems := make([]searchBySha256ResponseItem, len(audioFiles))
	for i, audioFile := range audioFiles {
		audioFilesResponseItems[i] = searchBySha256ResponseItem{
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

	c.JSON(http.StatusOK, searchBySha256Response{
		AudioFiles: audioFilesResponseItems,
	})
}
