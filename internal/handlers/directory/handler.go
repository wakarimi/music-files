package directory

import "music-files/internal/database/repository"

type DirHandler struct {
	DirRepo repository.DirRepositoryInterface
}

func NewDirHandler(dirRepo repository.DirRepositoryInterface) (h *DirHandler) {
	return &DirHandler{dirRepo}
}
