package audio_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/model/audio"
)

func (s Service) GetByDirAndName(tx *sqlx.Tx, dirID int, name string) (audio.Audio, error) {
	//TODO implement me
	panic("implement me")
}
