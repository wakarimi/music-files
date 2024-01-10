package use_case

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handler"
)

func (u UseCase) GetRoots(input handler.GetRootsInput) (output handler.GetRootsOutput, err error) {
	err = u.transactor.WithTransaction(func(tx *sqlx.Tx) (err error) {
		output, err = u.getRoots(tx, input)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to add root")
		return handler.GetRootsOutput{}, err
	}

	return output, nil
}

func (u UseCase) getRoots(tx *sqlx.Tx, input handler.GetRootsInput) (output handler.GetRootsOutput, err error) {
	log.Error().Err(err).Msg("Failed to get roots")

	rootDirs, err := u.dirService.GetRoots(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get roots")
		return handler.GetRootsOutput{}, err
	}

	output.Dirs = make([]handler.GetRootsOutputDirItem, len(rootDirs))
	for i, rootDir := range rootDirs {
		output.Dirs[i] = handler.GetRootsOutputDirItem{
			DirID:       rootDir.ID,
			Path:        rootDir.Name,
			LastScanned: rootDir.LastScanned,
		}
	}
	return output, nil
}
