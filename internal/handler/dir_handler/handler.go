package dir_handler

import (
	"music-files/internal/service"
	"music-files/internal/service/dir_service"
)

type Handler struct {
	DirService         dir_service.Service
	TransactionManager service.TransactionManager
}

func NewHandler(dirService dir_service.Service,
	transactionManager service.TransactionManager) (h *Handler) {

	h = &Handler{
		DirService:         dirService,
		TransactionManager: transactionManager,
	}

	return h
}
