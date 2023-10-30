package audio_file_handler

import (
	"music-files/internal/service"
	"music-files/internal/service/audio_file_service"
	"music-files/internal/service/file_processor_service"
)

type Handler struct {
	AudioFileService     audio_file_service.Service
	FileProcessorService file_processor_service.Service
	TransactionManager   service.TransactionManager
}

func NewHandler(audioFileService audio_file_service.Service,
	fileProcessorService file_processor_service.Service,
	transactionManager service.TransactionManager) (h *Handler) {

	h = &Handler{
		AudioFileService:     audioFileService,
		FileProcessorService: fileProcessorService,
		TransactionManager:   transactionManager,
	}

	return h
}
