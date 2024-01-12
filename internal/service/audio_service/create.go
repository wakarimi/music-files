package audio_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/model/audio"
)

func (s Service) Create(tx *sqlx.Tx, audioToCreate audio.Audio) (int, error) {
	//TODO implement me
	panic("implement me")
}
