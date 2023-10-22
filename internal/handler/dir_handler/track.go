package dir_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/handler/responses"
	"music-files/internal/models"
	"net/http"
)

// trackRequest is the request model for adding a new directory for tracking
type trackRequest struct {
	// Path to the directory on disk
	Path string `json:"path" bind:"required"`
}

// trackResponse is the response model after successfully adding a tracked directory
type trackResponse struct {
	// Unique identifier of the directory in the database
	DirId int `json:"dirId"`
	// Name of the directory
	Name string `json:"name"`
}

// Track
// @Summary Add a new tracked directory
// @Description Adds a new directory to the database for tracking
// @Tags Directories
// @Accept  json
// @Produce  json
// @Param   request body trackRequest true "Directory Data"
// @Success 201 {object} trackResponse
// @Failure 400 {object} responses.Error "Failed to decode request"
// @Failure 404 {object} responses.Error "Directory not found"
// @Failure 409 {object} responses.Error "Directory already tracked"
// @Failure 500 {object} responses.Error "Internal Server Error"
// @Router  /roots [post]
func (h *Handler) Track(c *gin.Context) {
	log.Debug().Msg("Adding a new directory tracking")

	var request trackRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error().Err(err)
		c.JSON(http.StatusBadRequest, responses.Error{
			Message: "Failed to encode request",
			Reason:  err.Error(),
		})
		return
	}

	var createdDir models.Directory
	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		dirToCreate := models.Directory{
			ParentDirId: nil,
			Name:        request.Path,
		}
		createdDir, err = h.DirService.Track(tx, dirToCreate)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to start tracking directory")
		if _, ok := err.(errors.NotFound); ok {
			c.JSON(http.StatusNotFound, responses.Error{
				Message: "Directory not found",
				Reason:  err.Error(),
			})
		} else if _, ok = err.(errors.Conflict); ok {
			c.JSON(http.StatusConflict, responses.Error{
				Message: "Directory already tracked",
				Reason:  err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, responses.Error{
				Message: "Failed to start tracking directory",
				Reason:  err.Error(),
			})
		}
		return
	}

	log.Debug().Msg("Directory added to tracked")
	c.JSON(http.StatusCreated, trackResponse{
		DirId: createdDir.DirId,
		Name:  createdDir.Name,
	})
}
