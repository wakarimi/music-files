package track

import "music-files/internal/database/repository"

type Handler struct {
	TrackRepo repository.TrackRepositoryInterface
}

func NewHandler(coverRepo repository.TrackRepositoryInterface) (h *Handler) {
	return &Handler{coverRepo}
}
