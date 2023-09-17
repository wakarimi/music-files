package cover

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handlers/types"
	"music-files/internal/models"
	"net/http"
	"strconv"
)

type readResponse struct {
	CoverId int    `json:"coverId"`
	Format  string `json:"format"`
	Size    int64  `json:"size"`
}

func (h *Handler) Read(c *gin.Context) {
	log.Debug().Msg("Fetching data about cover")

	coverIdStr := c.Param("coverId")
	coverId, err := strconv.Atoi(coverIdStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid coverId format")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error: "Invalid coverId format",
		})
		return
	}
	log.Debug().Int("coverId", coverId).Msg("Url parameter read successfully")

	var cover models.Cover

	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		cover, err = h.CoverService.Read(tx, coverId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch cover")
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error: "Failed to fetch cover",
		})
		return
	}

	log.Debug().Int("coverId", cover.CoverId).Str("relativePath", cover.RelativePath).Msg("Cover fetched successfully")
	c.JSON(http.StatusOK, readResponse{
		CoverId: cover.CoverId,
		Format:  cover.Format,
		Size:    cover.Size,
	})
}
