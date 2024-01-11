package use_case

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handler"
	"music-files/internal/internal_error"
	"music-files/internal/model/directory"
	"os"
)

func (u UseCase) ScanDir(input handler.ScanDirInput) (output handler.ScanDirOutput, err error) {
	log.Debug().Msg("Scanning directory")

	err = u.transactor.WithTransaction(func(tx *sqlx.Tx) (err error) {
		dirExists, err := u.dirService.IsExists(tx, input.DirID)
		if err != nil {
			return err
		}
		if !dirExists {
			err = internal_error.NotFound{EntityName: fmt.Sprintf("dir with id=%d", input.DirID)}
			return err
		}

		absolutePath, err := u.dirService.CalcAbsolutePath(tx, input.DirID)
		if err != nil {
			return err
		}
		dirExistsOnDisk, err := u.dirService.IsExistsOnDisk(absolutePath)
		if err != nil {
			return err
		}
		if !dirExistsOnDisk {
			err = internal_error.NotFound{EntityName: fmt.Sprintf("dir on disk with path=%s", absolutePath)}
			return err
		}

		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to scan dir")
		return handler.ScanDirOutput{}, err
	}

	go u.scanDir(input.DirID)

	return handler.ScanDirOutput{}, nil
}

func (u UseCase) scanDir(dirID int) {
	var subDirs []directory.Directory
	err := u.transactor.WithTransaction(func(tx *sqlx.Tx) (err error) {
		err = u.scanDirActualizeSubDirs(tx, dirID)
		if err != nil {
			return err
		}
		subDirs, err = u.scanDirReadSubDirs(tx, dirID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return
	}

	for _, subDir := range subDirs {
		u.scanDir(subDir.ID)
	}

	err = u.transactor.WithTransaction(func(tx *sqlx.Tx) (err error) {
		err = u.scanDirScanContent(tx, dirID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return
	}
}

func (u UseCase) scanDirExists(tx *sqlx.Tx, dirID int) (exists bool, err error) {
	dirExistsInDatabase, err := u.dirService.IsExists(tx, dirID)
	if err != nil {
		log.Error().Err(err).Int("dirID", dirID).Msg("Failed to check directory existence in database")
		return false, err
	}
	if !dirExistsInDatabase {
		log.Error().Err(err).Int("dirID", dirID).Msg("Directory not found in database")
		return false, err
	}

	absolutePath, err := u.dirService.CalcAbsolutePath(tx, dirID)
	if err != nil {
		log.Error().Err(err).Int("dirID", dirID).Msg("Failed to calculate absolute path")
		return false, err
	}
	existOnDisk, err := u.dirService.IsExistsOnDisk(absolutePath)
	if err != nil {
		log.Error().Err(err).Int("dirID", dirID).Msg("Failed to check directory existence on disk")
		return false, err
	}
	if !existOnDisk {
		log.Error().Err(err).Int("dirID", dirID).Msg("Directory not found on disk")
		return false, nil
	}

	return true, nil
}

func (u UseCase) scanDirActualizeSubDirs(tx *sqlx.Tx, dirID int) error {
	absolutePath, err := u.dirService.CalcAbsolutePath(tx, dirID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to calculate absolute path")
		return err
	}

	entries, err := os.ReadDir(absolutePath)
	if err != nil {
		log.Error().Err(err).Msg("failed to read directory from disk")
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			alreadyInDatabase, err := u.dirService.IsExistsByParentAndName(tx, &dirID, entry.Name())
			if err != nil {
				log.Error().Err(err).Msg("Failed to check directory existence")
				return err
			}
			if !alreadyInDatabase {
				_, err = u.dirService.Create(tx, directory.Directory{
					ParentDirID: &dirID,
					Name:        entry.Name(),
				})
				if err != nil {
					log.Error().Err(err).Msg("Failed to create directory in database")
					return err
				}
			}
		}
	}

	subDirs, err := u.dirService.GetSubDirs(tx, dirID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read subdirectories")
		return err
	}
	for _, subDir := range subDirs {
		subDirFoundOnDisk := false

		for _, entry := range entries {
			if entry.IsDir() {
				if subDir.Name == entry.Name() {
					subDirFoundOnDisk = true
				}
			}
		}

		if !subDirFoundOnDisk {
			err = u.deleteRootDeleteDirRecursive(tx, subDir.ID)
			if err != nil {
				log.Error().Err(err).Msg("Failed to delete directory from disk")
				return err
			}
		}
	}

	return nil
}

func (u UseCase) scanDirDeleteDirRecursive(tx *sqlx.Tx, dirID int) (err error) {
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

func (u UseCase) scanDirReadSubDirs(tx *sqlx.Tx, dirID int) ([]directory.Directory, error) {

}

func (u UseCase) scanDirScanContent(tx *sqlx.Tx, dirID int) error {

}
