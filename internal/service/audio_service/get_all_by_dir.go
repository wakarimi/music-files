package audio_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/model/audio"
)

func (s Service) GetAllByDir(tx *sqlx.Tx, dirID int) ([]audio.Audio, error) {
	//TODO implement me
	panic("implement me")
}
