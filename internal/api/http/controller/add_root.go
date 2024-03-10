package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type addRootRequest struct {
	Path string `json:"path" binding:"required"`
}

func (c *Controller) AddRoot(ctx *gin.Context) {
	c.log.Debug().Msg("Creating new root directory")

	var request addRootRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		c.log.Error().Err(err).Msg("Invalid request body")
		messageID := "InvalidRequestBody"
		ctx.JSON(http.StatusBadRequest, c.ConstructError(ctx, messageID, err))
		return
	}
}
