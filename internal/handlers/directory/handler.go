package directory

import "music-files/internal/database/repository"

type DirHandler struct {
	DirRepo   repository.DirRepositoryInterface
	CoverRepo repository.CoverRepositoryInterface
	TrackRepo repository.TrackRepositoryInterface
}

func NewDirHandler(dirRepo repository.DirRepositoryInterface,
	coverRepo repository.CoverRepositoryInterface,
	trackRepo repository.TrackRepositoryInterface) (h *DirHandler) {

	return &DirHandler{dirRepo, coverRepo, trackRepo}
}
