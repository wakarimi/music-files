package directory

import "music-files/internal/database/repository"

type Handler struct {
	DirRepo   repository.DirRepositoryInterface
	CoverRepo repository.CoverRepositoryInterface
	TrackRepo repository.TrackRepositoryInterface
}

func NewHandler(dirRepo repository.DirRepositoryInterface,
	coverRepo repository.CoverRepositoryInterface,
	trackRepo repository.TrackRepositoryInterface) (h *Handler) {

	return &Handler{dirRepo, coverRepo, trackRepo}
}
