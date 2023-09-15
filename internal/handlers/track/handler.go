package track

import (
	"music-files/internal/service"
	"music-files/internal/service/dir_service"
	"music-files/internal/service/track_service"
)

type Handler struct {
	TransactionManager service.TransactionManager
	TrackService       track_service.Service
	DirService         dir_service.Service
}

func NewHandler(transactionManager service.TransactionManager,
	trackService track_service.Service,
	dirService dir_service.Service) (h *Handler) {

	h = &Handler{
		TransactionManager: transactionManager,
		TrackService:       trackService,
		DirService:         dirService,
	}

	return h
}
