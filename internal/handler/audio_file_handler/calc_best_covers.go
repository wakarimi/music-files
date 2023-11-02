package audio_file_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handler/response"
	"net/http"
)

// calcBestCoversRequest are the input parameters for which you need to find the top covers
type calcBestCoversRequest struct {
	// List of audio files to find their most suitable covers
	AudioFiles []int `json:"audioFiles" bind:"required"`
}

// calcBestCoversResponse is the calculated top covers
type calcBestCoversResponse struct {
	// Top suitable covers by frequency of occurrence
	Covers []int `json:"coversTop"`
}

// CalcBestCovers retrieves top of covers for audio file
// @Summary Retrieve top of covers for audio file
// @Description Retrieves a top of covers for audio file
// @Tags Covers
// @Accept  json
// @Produce  json
// @Param   request body calcBestCoversRequest true "Directory Data"
// @Success 200 {object} calcBestCoversResponse
// @Failure 500 {object} response.Error "Internal Server Error"
// @Router /audio-files/covers-top [put]
func (h *Handler) CalcBestCovers(c *gin.Context) {
	log.Debug().Msg("Getting cover for audioFile")

	var request calcBestCoversRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error().Err(err).Msg("Failed to encode request")
		c.JSON(http.StatusBadRequest, response.Error{
			Message: "Failed to encode request",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Msg("Request encoded successfully")

	audioFileIds := make([]int, len(request.AudioFiles))
	for i, requestItem := range request.AudioFiles {
		audioFileIds[i] = requestItem
	}
	var covers []int
	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		covers, err = h.FileProcessorService.CalcBestCovers(tx, audioFileIds)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get cover")
		c.JSON(http.StatusInternalServerError, response.Error{
			Message: "Failed to get cover",
			Reason:  err.Error(),
		})
		return
	}

	log.Debug().Msg("Best covers calculated successfully")
	c.JSON(http.StatusOK, calcBestCoversResponse{
		Covers: covers,
	})
}
