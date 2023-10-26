package song_handler

import (
	"music-files/internal/service"
	"music-files/internal/service/file_processor_service"
	"music-files/internal/service/song_service"
)

type Handler struct {
	SongService          song_service.Service
	FileProcessorService file_processor_service.Service
	TransactionManager   service.TransactionManager
}

func NewHandler(songService song_service.Service,
	fileProcessorService file_processor_service.Service,
	transactionManager service.TransactionManager) (h *Handler) {

	h = &Handler{
		SongService:          songService,
		FileProcessorService: fileProcessorService,
		TransactionManager:   transactionManager,
	}

	return h
}
