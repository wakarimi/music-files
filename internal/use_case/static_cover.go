package use_case

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handler"
	"music-files/internal/internal_error"
	"path/filepath"
)

func (u UseCase) StaticCover(input handler.StaticCoverInput) (output handler.StaticCoverOutput, err error) {
	err = u.transactor.WithTransaction(func(tx *sqlx.Tx) (err error) {
		output, err = u.staticCover(tx, input)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get static cover")
		return handler.StaticCoverOutput{}, err
	}

	return output, nil
}

func (u UseCase) staticCover(tx *sqlx.Tx, input handler.StaticCoverInput) (handler.StaticCoverOutput, error) {
	log.Debug().Msg("Getting static cover")

	existsInDatabase, err := u.coverService.IsExists(tx, input.CoverID)
	if err != nil {
		log.Error().Err(err).Int("coverId", input.CoverID).Msg("Failed to check cover existence")
		return handler.StaticCoverOutput{}, err
	}
	if !existsInDatabase {
		err := internal_error.NotFound{fmt.Sprintf("cover with id=%d", input.CoverID)}
		log.Error().Err(err).Int("coverId", input.CoverID).Msg("Cover not found")
		return handler.StaticCoverOutput{}, err
	}

	coverModel, err := u.coverService.Get(tx, input.CoverID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get cover")
		return handler.StaticCoverOutput{}, err
	}

	absolutePathToDir, err := u.dirService.CalcAbsolutePath(tx, coverModel.DirID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to calculate absolute path")
		return handler.StaticCoverOutput{}, err
	}
	absolutePath := filepath.Join(absolutePathToDir, coverModel.Filename)

	existsOnDisk, err := u.coverService.IsExistsOnDisk(absolutePath)
	if err != nil {
		log.Error().Err(err).Int("coverId", input.CoverID).Msg("Failed to check cover existence on disk")
		return handler.StaticCoverOutput{}, err
	}
	if !existsOnDisk {
		log.Error().Err(err).Int("coverId", input.CoverID).Msg("Cover not found on disk")
		return handler.StaticCoverOutput{}, err
	}

	mime, err := u.coverService.GetMimeValue(absolutePath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get mime type")
		return handler.StaticCoverOutput{}, err
	}

	log.Debug().Msg("Static link to cover got")
	return handler.StaticCoverOutput{
		AbsolutePath: absolutePath,
		Mime:         mime,
	}, nil
}
