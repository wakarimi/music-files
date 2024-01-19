package handler

import "github.com/gin-gonic/gin"

type GetCoversRatingResponse struct {
	CoversRating []int `json:"coversRating"`
}

func (h Handler) GetCoversRating(c *gin.Context) {

}
