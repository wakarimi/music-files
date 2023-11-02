package dir_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/handler/response"
	"music-files/internal/model"
	"net/http"
	"strconv"
	"time"
)

// getDirResponse is the model for each item in the getRoots response
type getDirResponse struct {
	// Unique identifier for the directory
	DirId int `json:"dirId"`
	// Name of the directory
	Name string `json:"name"`
	// Absolute path to directory
	AbsolutePath string `json:"absolutePath"`
	// Last time the directory was scanned
	LastScanned *time.Time `json:"lastScanned,omitempty"`
}

// GetDir
// @Summary Retrieve a directory by ID
// @Description Retrieves detailed information about a directory by its ID
// @Tags Directories
// @Accept  json
// @Produce  json
// @Param   dirId     path    int     true        "Dir ID"
// @Success 200 {object} getDirResponse
// @Failure 400 {object} response.Error "Invalid dirId format"
// @Failure 404 {object} response.Error "Directory not found"
// @Failure 500 {object} response.Error "Internal Server Error"
// @Router /dirs/{dirId} [get]
func (h *Handler) GetDir(c *gin.Context) {
	log.Debug().Msg("Getting directory")

	dirIdStr := c.Param("dirId")
	dirId, err := strconv.Atoi(dirIdStr)
	if err != nil {
		log.Error().Err(err).Str("dirIdStr", dirIdStr).Msg("Invalid dirId format")
		c.JSON(http.StatusBadRequest, response.Error{
			Message: "Invalid dirId format",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("dirId", dirId).Msg("Url parameter read successfully")

	var dir model.Directory
	var absolutePath string
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		dir, err = h.DirService.GetDir(tx, dirId)
		if err != nil {
			return err
		}
		absolutePath, err = h.DirService.AbsolutePath(tx, dirId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get directory")
		if _, ok := err.(errors.NotFound); ok {
			c.JSON(http.StatusNotFound, response.Error{
				Message: "Directory not found",
				Reason:  err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, response.Error{
				Message: "Failed to get content",
				Reason:  err.Error(),
			})
		}
		return
	}

	log.Debug().Msg("Directory content got successfully")
	c.JSON(http.StatusOK, getDirResponse{
		DirId:        dir.DirId,
		Name:         dir.Name,
		AbsolutePath: absolutePath,
		LastScanned:  dir.LastScanned,
	})
}
