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

type trackRequest struct {
	Path string `json:"path" bind:"required"`
}

type trackResponse struct {
	DirId int    `db:"dirId"`
	Name  string `db:"name"`
}

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
		log.Error().Err(err).Msg("Failed to start tracking directory")
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
