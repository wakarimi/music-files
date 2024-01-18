package use_case

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/internal_error"
	"music-files/internal/model/directory"
	"strings"
)

type AddRootInput struct {
	Path string
}

type AddRootOutput struct {
	DirID int
	Path  string
}

func (u UseCase) AddRoot(input AddRootInput) (output AddRootOutput, err error) {
	err = u.transactor.WithTransaction(func(tx *sqlx.Tx) (err error) {
		output, err = u.addRoot(tx, input)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to add root")
		return AddRootOutput{}, err
	}

	return output, nil
}

func (u UseCase) addRoot(tx *sqlx.Tx, input AddRootInput) (output AddRootOutput, err error) {
	log.Debug().Msg("Adding root directory")

	for strings.HasSuffix(input.Path, "/") {
		input.Path = strings.TrimSuffix(input.Path, "/")
	}

	alreadyTracked, err := u.dirService.IsTracked(tx, input.Path)
	if err != nil {
		log.Error().Err(err).Msg("Couldn't check if the directory tracked")
		return AddRootOutput{}, err
	}
	if alreadyTracked {
		err := internal_error.Conflict{Message: "directory already tracked"}
		log.Error().Err(err).Msg("Directory already tracked")
		return AddRootOutput{}, err
	}

	dirExistsOnDisk, err := u.dirService.IsExistsOnDisk(input.Path)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check directory on disk")
		return AddRootOutput{}, err
	}
	if !dirExistsOnDisk {
		err := internal_error.NotFound{EntityName: fmt.Sprintf("dir with path %s", input.Path)}
		log.Error().Err(err).Msg("Directory not found")
		return AddRootOutput{}, err
	}

	containedRoots, err := u.dirService.ContainedRoots(tx, input.Path)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get contained roots")
		return AddRootOutput{}, err
	}

	dirToCreate := directory.Directory{
		Name: input.Path,
	}
	createdDirID, err := u.dirService.Create(tx, dirToCreate)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create dir")
		return AddRootOutput{}, err
	}
	createdDir, err := u.dirService.Get(tx, createdDirID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get created dir")
		return AddRootOutput{}, err
	}

	for _, containedRoot := range containedRoots {
		err = u.dirService.MergeRoots(tx, createdDirID, containedRoot.ID)
		if err != nil {
			return AddRootOutput{}, err
		}
	}

	log.Debug().Msg("Root dir added")
	return AddRootOutput{
		DirID: createdDir.ID,
		Path:  createdDir.Name,
	}, nil
}
