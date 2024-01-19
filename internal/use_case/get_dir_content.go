package use_case

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/internal_error"
	"time"
)

type GetDirContentInput struct {
	DirID int
}

type GetDirContentOutputDirs struct {
	ID          int
	Name        string
	LastScanned *time.Time
}

type GetDirContentOutputAudios struct {
	ID                int
	SHA256            string
	LastContentUpdate time.Time
}

type GetDirContentOutput struct {
	Dirs   []GetDirContentOutputDirs
	Audios []GetDirContentOutputAudios
}

func (u UseCase) GetDirContent(input GetDirContentInput) (output GetDirContentOutput, err error) {
	err = u.transactor.WithTransaction(func(tx *sqlx.Tx) (err error) {
		output, err = u.getDirContent(tx, input)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to add root")
		return GetDirContentOutput{}, err
	}

	return output, nil
}
func (u UseCase) getDirContent(tx *sqlx.Tx, input GetDirContentInput) (GetDirContentOutput, error) {
	log.Debug().Msg("Getting dir's content")

	exists, err := u.dirService.IsExists(tx, input.DirID)
	if err != nil {
		log.Error().Err(err).Int("dirId", input.DirID).Msg("Failed to check dir existence")
		return GetDirContentOutput{}, err
	}
	if !exists {
		err := internal_error.NotFound{fmt.Sprintf("directory with id=%d", input.DirID)}
		log.Error().Err(err).Int("dirId", input.DirID).Msg("Directory not found")
		return GetDirContentOutput{}, err
	}

	dirs, err := u.dirService.GetSubDirs(tx, input.DirID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get subdirectories")
		return GetDirContentOutput{}, err
	}

	audios, err := u.audioService.GetAllByDir(tx, input.DirID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get dir's audios")
		return GetDirContentOutput{}, err
	}

	dirsResponse := make([]GetDirContentOutputDirs, len(dirs))
	for i, dir := range dirs {
		dirsResponse[i] = GetDirContentOutputDirs{
			ID:          dir.ID,
			Name:        dir.Name,
			LastScanned: dir.LastScanned,
		}
	}

	audiosResponse := make([]GetDirContentOutputAudios, len(audios))
	for i, aud := range audios {
		audiosResponse[i] = GetDirContentOutputAudios{
			ID:                aud.ID,
			SHA256:            aud.SHA256,
			LastContentUpdate: aud.LastContentUpdate,
		}
	}

	return GetDirContentOutput{
		Dirs:   dirsResponse,
		Audios: audiosResponse,
	}, nil
}
