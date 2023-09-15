package directory

import (
	"music-files/internal/service"
	"music-files/internal/service/cover_service"
	"music-files/internal/service/dir_service"
	"music-files/internal/service/track_service"
)

type Handler struct {
	TransactionManager service.TransactionManager
	DirService         dir_service.Service
	CoverService       cover_service.Service
	TrackService       track_service.Service
}

func NewHandler(transactionManager service.TransactionManager,
	dirService dir_service.Service,
	coverService cover_service.Service,
	trackService track_service.Service) (h *Handler) {

	h = &Handler{
		TransactionManager: transactionManager,
		DirService:         dirService,
		CoverService:       coverService,
		TrackService:       trackService,
	}
	return h
}
