package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"music-files/internal/database/repository"
	"music-files/internal/handlers/types"
	"net/http"
	"strconv"
	"time"
)

type CoverGetResponse struct {
	CoverId   int       `json:"coverId"`
	Size      int64     `json:"size"`
	Format    string    `json:"format"`
	DateAdded time.Time `json:"dateAdded"`
}

func CoverGet(c *gin.Context) {
	coverIdStr := c.Param("coverId")

	coverId, err := strconv.Atoi(coverIdStr)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, types.Error{
			Error: "Invalid coverId format",
		})
		return
	}

	cover, err := repository.GetCoverById(coverId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to get cover",
		})
		return
	}

	c.JSON(http.StatusOK, CoverGetResponse{
		CoverId:   cover.CoverId,
		Size:      cover.Size,
		Format:    cover.Format,
		DateAdded: cover.DateAdded,
	})

}
