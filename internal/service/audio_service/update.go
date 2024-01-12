package audio_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/model/audio"
)

func (s Service) Update(tx *sqlx.Tx, audioID int, audioToUpdate audio.Audio) error {
	//TODO implement me
	panic("implement me")
}
