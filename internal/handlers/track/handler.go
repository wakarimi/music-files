package track

import "music-files/internal/database/repository"

type Handler struct {
	TrackRepo repository.TrackRepositoryInterface
	DirRepo   repository.DirRepositoryInterface
}

func NewHandler(coverRepo repository.TrackRepositoryInterface, dirRepo repository.DirRepositoryInterface) (h *Handler) {
	return &Handler{coverRepo, dirRepo}
}
