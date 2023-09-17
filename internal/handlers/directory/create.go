package directory

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handlers/types"
	"music-files/internal/models"
	"net/http"
)

// createRequest godoc
// @Description Request structure to create a new directory
// @Property Path (string) The path to the directory
type createRequest struct {
	Path string `json:"path" bind:"required"`
}

// Create godoc
// @Summary Create a directory
// @Description Adds a new directory
// @Tags Directories
// @Accept  json
// @Produce  json
// @Param body body createRequest true "Details for the new directory"
// @Success 201 {string} none "Successfully created directory"
// @Failure 400 {object} types.ErrorResponse "Failed to encode request or Invalid input"
// @Failure 500 {object} types.ErrorResponse "Failed to create directory"
// @Router /dirs [post]
func (h *Handler) Create(c *gin.Context) {
	log.Debug().Msg("Creating new directory")

	var request createRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error().Err(err).Msg("Failed to encode request")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error: "Failed to encode request",
		})
		return
	}

	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		dir := models.Directory{
			Path: request.Path,
		}
		err = h.DirService.Create(tx, dir)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create directory")
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error: "Failed to create directory",
		})
		return
	}

	log.Debug().Msg("Directory added successfully")
	c.Status(http.StatusCreated)
}
