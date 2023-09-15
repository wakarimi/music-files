package cover

import (
	"music-files/internal/service"
	"music-files/internal/service/cover_service"
	"music-files/internal/service/dir_service"
)

type Handler struct {
	TransactionManager service.TransactionManager
	CoverService       cover_service.Service
	DirService         dir_service.Service
}

func NewHandler(transactionManager service.TransactionManager,
	coverService cover_service.Service,
	dirService dir_service.Service) (h *Handler) {

	h = &Handler{
		TransactionManager: transactionManager,
		CoverService:       coverService,
		DirService:         dirService,
	}

	return h
}
