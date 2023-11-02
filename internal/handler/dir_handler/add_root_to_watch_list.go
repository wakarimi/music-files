package dir_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/handler/response"
	"music-files/internal/model"
	"net/http"
)

// addRootToWatchListRequest is the request model for adding a new directory for tracking
type addRootToWatchListRequest struct {
	// Path to the directory on disk
	Path string `json:"path" bind:"required"`
}

// addRootToWatchListResponse is the response model after successfully adding a tracked directory
type addRootToWatchListResponse struct {
	// Unique identifier of the directory in the database
	DirId int `json:"dirId"`
	// Name of the directory
	Name string `json:"name"`
}

// AddRootToWatchList
// @Summary Add a new tracked directory
// @Description Adds a new directory to the database for tracking
// @Tags Directories
// @Accept  json
// @Produce  json
// @Param   request body addRootToWatchListRequest true "Directory Data"
// @Success 201 {object} addRootToWatchListResponse
// @Failure 400 {object} response.Error "Failed to decode request"
// @Failure 404 {object} response.Error "Directory not found"
// @Failure 409 {object} response.Error "Directory already tracked"
// @Failure 500 {object} response.Error "Internal Server Error"
// @Router  /roots [post]
func (h *Handler) AddRootToWatchList(c *gin.Context) {
	log.Debug().Msg("Adding a root directory to the watch list")

	var request addRootToWatchListRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error().Err(err).Msg("Failed to encode request")
		c.JSON(http.StatusBadRequest, response.Error{
			Message: "Failed to encode request",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Str("path", request.Path).Msg("Request encoded successfully")

	var createdDir model.Directory
	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		dirToCreate := model.Directory{
			ParentDirId: nil,
			Name:        request.Path,
		}
		createdDir, err = h.DirService.AddRootToWatchList(tx, dirToCreate)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to add directory to watch list")
		if _, ok := err.(errors.NotFound); ok {
			c.JSON(http.StatusNotFound, response.Error{
				Message: "Directory not found on disk",
				Reason:  err.Error(),
			})
		} else if _, ok = err.(errors.Conflict); ok {
			c.JSON(http.StatusConflict, response.Error{
				Message: "The directory is already being tracked",
				Reason:  err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, response.Error{
				Message: "Failed to add directory to watch list",
				Reason:  err.Error(),
			})
		}
		return
	}

	log.Debug().Msg("Directory added to watch list successfully")
	c.JSON(http.StatusCreated, addRootToWatchListResponse{
		DirId: createdDir.DirId,
		Name:  createdDir.Name,
	})
}
