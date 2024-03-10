package controller

import (
	"github.com/gin-gonic/gin"
	"music-files/pkg/httpresp"
)

func (c *Controller) ConstructError(ctx *gin.Context, messageID string, err error) *httpresp.Error {
	lang := ctx.MustGet("lang").(string)

	return &httpresp.Error{
		Message:          c.localizer.English(messageID),
		LocalizedMessage: c.localizer.TryLocalize(lang, messageID),
		Error:            err.Error(),
	}
}
