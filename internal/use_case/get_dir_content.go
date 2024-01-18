package use_case

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handler"
	"music-files/internal/internal_error"
)

func (u UseCase) GetDirContent(input handler.GetDirContentInput) (output handler.GetDirContentOutput, err error) {
	err = u.transactor.WithTransaction(func(tx *sqlx.Tx) (err error) {
		output, err = u.getDirContent(tx, input)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to add root")
		return handler.GetDirContentOutput{}, err
	}

	return output, nil
}
func (u UseCase) getDirContent(tx *sqlx.Tx, input handler.GetDirContentInput) (handler.GetDirContentOutput, error) {
	log.Debug().Msg("Getting dir's content")

	exists, err := u.dirService.IsExists(tx, input.DirID)
	if err != nil {
		log.Error().Err(err).Int("dirId", input.DirID).Msg("Failed to check dir existence")
		return handler.GetDirContentOutput{}, err
	}
	if !exists {
		err := internal_error.NotFound{fmt.Sprintf("directory with id=%d", input.DirID)}
		log.Error().Err(err).Int("dirId", input.DirID).Msg("Directory not found")
		return handler.GetDirContentOutput{}, err
	}

	dirs, err := u.dirService.GetSubDirs(tx, input.DirID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get subdirectories")
		return handler.GetDirContentOutput{}, err
	}

	audios, err := u.audioService.GetAllByDir(tx, input.DirID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get dir's audios")
		return handler.GetDirContentOutput{}, err
	}

	dirsResponse := make([]handler.GetDirContentOutputDirs, len(dirs))
	for i, dir := range dirs {
		dirsResponse[i] = handler.GetDirContentOutputDirs{
			ID:          dir.ID,
			Name:        dir.Name,
			LastScanned: dir.LastScanned,
		}
	}

	audiosResponse := make([]handler.GetDirContentOutputAudios, len(audios))
	for i, aud := range audios {
		audiosResponse[i] = handler.GetDirContentOutputAudios{
			ID:                aud.ID,
			DirID:             aud.DirID,
			DurationMs:        aud.DurationMs,
			SHA256:            aud.SHA256,
			LastContentUpdate: aud.LastContentUpdate,
		}
	}

	return handler.GetDirContentOutput{
		Dirs:   dirsResponse,
		Audios: audiosResponse,
	}, nil
}
