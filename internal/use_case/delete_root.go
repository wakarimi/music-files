package use_case

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handler"
	"music-files/internal/internal_error"
)

func (u UseCase) DeleteRoot(input handler.DeleteRootInput) (output handler.DeleteRootOutput, err error) {
	err = u.transactor.WithTransaction(func(tx *sqlx.Tx) (err error) {
		output, err = u.deleteRoot(tx, input)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to add root")
		return handler.DeleteRootOutput{}, err
	}

	return output, err
}

func (u UseCase) deleteRoot(tx *sqlx.Tx, input handler.DeleteRootInput) (output handler.DeleteRootOutput, err error) {
	log.Debug().Msg("Deleting root directory")

	exists, err := u.dirService.IsExists(tx, input.DirID)
	if err != nil {
		log.Error().Err(err).Int("dirId", input.DirID).Msg("Failed to check dir existence")
		return handler.DeleteRootOutput{}, err
	}
	if !exists {
		err := internal_error.NotFound{fmt.Sprintf("directory with id=%d", input.DirID)}
		log.Error().Err(err).Int("dirId", input.DirID).Msg("Directory not found")
		return handler.DeleteRootOutput{}, err
	}

	isRoot, err := u.dirService.IsRoot(tx, input.DirID)
	if err != nil {
		log.Error().Err(err).Int("dirId", input.DirID).Msg("Could not check whether the directory is the root")
		return handler.DeleteRootOutput{}, err
	}
	if !isRoot {
		err := internal_error.BadRequest{fmt.Sprintf("directory with id=%d is not root", input.DirID)}
		log.Error().Err(err).Int("dirId", input.DirID).Msg("Directory is not root")
		return handler.DeleteRootOutput{}, err
	}

	err = u.deleteRootDeleteDirRecursive(tx, input.DirID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete dir recursive")
		return handler.DeleteRootOutput{}, err
	}

	return output, err
}

func (u UseCase) deleteRootDeleteDirRecursive(tx *sqlx.Tx, dirID int) (err error) {
	log.Debug().Int("dirId", dirID).Msg("Deleting directory")

	err = u.audioService.DeleteAllByDir(tx, dirID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete directory's audios")
		return err
	}

	err = u.coverService.DeleteAllByDir(tx, dirID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete directory's covers")
		return err
	}

	subDirs, err := u.dirService.GetSubDirs(tx, dirID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get subdirs")
		return err
	}

	for _, subDir := range subDirs {
		err := u.deleteRootDeleteDirRecursive(tx, subDir.ID)
		if err != nil {
			log.Error().Err(err).Int("subDirId", subDir.ID).Msg("Failed to delete subdirectory")
			return err
		}
	}

	err = u.dirService.Delete(tx, dirID)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirID).Msg("Failed to delete dir")
		return err
	}

	return nil
}
