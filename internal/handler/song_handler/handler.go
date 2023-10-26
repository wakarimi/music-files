package song_handler

import (
	"music-files/internal/service"
	"music-files/internal/service/song_service"
)

type Handler struct {
	SongService        song_service.Service
	TransactionManager service.TransactionManager
}

func NewHandler(songService song_service.Service,
	transactionManager service.TransactionManager) (h *Handler) {

	h = &Handler{
		SongService:        songService,
		TransactionManager: transactionManager,
	}

	return h
}
