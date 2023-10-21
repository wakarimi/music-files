package song_service

import (
	"music-files/internal/database/repository/song_repo"
)

type Service struct {
	SongRepo song_repo.Repo
}

func NewService(songRepo song_repo.Repo) (s *Service) {

	s = &Service{
		SongRepo: songRepo,
	}

	return s
}
