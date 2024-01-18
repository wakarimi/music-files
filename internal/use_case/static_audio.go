package use_case

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handler"
	"music-files/internal/internal_error"
	"path/filepath"
)

func (u UseCase) StaticAudio(input handler.StaticAudioInput) (output handler.StaticAudioOutput, err error) {
	err = u.transactor.WithTransaction(func(tx *sqlx.Tx) (err error) {
		output, err = u.staticAudio(tx, input)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get static audio")
		return handler.StaticAudioOutput{}, err
	}

	return output, nil
}

func (u UseCase) staticAudio(tx *sqlx.Tx, input handler.StaticAudioInput) (handler.StaticAudioOutput, error) {
	log.Debug().Msg("Getting static audio")

	existsInDatabase, err := u.audioService.IsExists(tx, input.AudioID)
	if err != nil {
		log.Error().Err(err).Int("audioId", input.AudioID).Msg("Failed to check audio existence")
		return handler.StaticAudioOutput{}, err
	}
	if !existsInDatabase {
		err := internal_error.NotFound{fmt.Sprintf("audio with id=%d", input.AudioID)}
		log.Error().Err(err).Int("audioId", input.AudioID).Msg("Audio not found")
		return handler.StaticAudioOutput{}, err
	}

	audioModel, err := u.audioService.Get(tx, input.AudioID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get audio")
		return handler.StaticAudioOutput{}, err
	}

	absolutePathToDir, err := u.dirService.CalcAbsolutePath(tx, audioModel.DirID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to calculate absolute path")
		return handler.StaticAudioOutput{}, err
	}
	absolutePath := filepath.Join(absolutePathToDir, audioModel.Filename)

	existsOnDisk, err := u.audioService.IsExistsOnDisk(absolutePath)
	if err != nil {
		log.Error().Err(err).Int("audioId", input.AudioID).Msg("Failed to check audio existence on disk")
		return handler.StaticAudioOutput{}, err
	}
	if !existsOnDisk {
		log.Error().Err(err).Int("audioId", input.AudioID).Msg("Audio not found on disk")
		return handler.StaticAudioOutput{}, err
	}

	mime, err := u.audioService.GetMimeValue(absolutePath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get mime type")
		return handler.StaticAudioOutput{}, err
	}

	log.Debug().Msg("Static link to audio got")
	return handler.StaticAudioOutput{
		AbsolutePath: absolutePath,
		Mime:         mime,
	}, nil
}
