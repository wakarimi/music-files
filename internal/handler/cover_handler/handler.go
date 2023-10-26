package cover_handler

import (
	"music-files/internal/service"
	"music-files/internal/service/cover_service"
	"music-files/internal/service/file_processor_service"
)

type Handler struct {
	CoverService         cover_service.Service
	FileProcessorService file_processor_service.Service
	TransactionManager   service.TransactionManager
}

func NewHandler(coverService cover_service.Service,
	fileProcessorService file_processor_service.Service,
	transactionManager service.TransactionManager) (h *Handler) {

	h = &Handler{
		CoverService:         coverService,
		FileProcessorService: fileProcessorService,
		TransactionManager:   transactionManager,
	}

	return h
}
