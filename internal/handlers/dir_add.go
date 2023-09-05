package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"music-files/internal/database/repository"
	"music-files/internal/handlers/types"
	"net/http"
)

type DirAddRequest struct {
	Path string `json:"path" binding:"required"`
}

func DirAdd(c *gin.Context) {
	var request DirAddRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, types.Error{
			Error: err.Error(),
		})
		return
	}

	exists, err := repository.DirExist(request.Path)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to check directory",
		})
		return
	}
	if exists {
		log.Println(err)
		c.JSON(http.StatusConflict, types.Error{
			Error: "Directory already added",
		})
		return
	}

	_, err = repository.InsertDir(request.Path)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to add directory",
		})
		return
	}

	c.Status(http.StatusCreated)
}
