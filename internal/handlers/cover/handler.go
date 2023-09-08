package cover

import "music-files/internal/database/repository"

type Handler struct {
	CoverRepo repository.CoverRepositoryInterface
}

func NewHandler(coverRepo repository.CoverRepositoryInterface) (h *Handler) {
	return &Handler{coverRepo}
}
