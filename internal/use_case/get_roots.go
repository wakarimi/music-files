package use_case

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"time"
)

type GetRootsInput struct {
}

type GetRootsOutputDirItem struct {
	DirID       int
	Path        string
	LastScanned *time.Time
}

type GetRootsOutput struct {
	Dirs []GetRootsOutputDirItem
}

func (u UseCase) GetRoots(input GetRootsInput) (output GetRootsOutput, err error) {
	err = u.transactor.WithTransaction(func(tx *sqlx.Tx) (err error) {
		output, err = u.getRoots(tx, input)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to add root")
		return GetRootsOutput{}, err
	}

	return output, nil
}

func (u UseCase) getRoots(tx *sqlx.Tx, input GetRootsInput) (output GetRootsOutput, err error) {
	log.Error().Err(err).Msg("Failed to get roots")

	rootDirs, err := u.dirService.GetRoots(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get roots")
		return GetRootsOutput{}, err
	}

	output.Dirs = make([]GetRootsOutputDirItem, len(rootDirs))
	for i, rootDir := range rootDirs {
		output.Dirs[i] = GetRootsOutputDirItem{
			DirID:       rootDir.ID,
			Path:        rootDir.Name,
			LastScanned: rootDir.LastScanned,
		}
	}
	return output, nil
}
