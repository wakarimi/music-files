package cover

import "music-files/internal/database/repository"

type Handler struct {
	CoverRepo repository.CoverRepositoryInterface
	DirRepo   repository.DirRepositoryInterface
}

func NewHandler(coverRepo repository.CoverRepositoryInterface, dirRepo repository.DirRepositoryInterface) (h *Handler) {
	return &Handler{coverRepo, dirRepo}
}
