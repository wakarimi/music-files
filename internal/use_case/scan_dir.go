package use_case

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handler"
	"music-files/internal/internal_error"
	"music-files/internal/model/directory"
	"os"
	"path/filepath"
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
		subDirs, err = u.dirService.GetSubDirs(tx, dirID)
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
		log.Error().Err(err).Int("dirId", dirID).Msg("Failed to check directory existence in database")
		return false, err
	}
	if !dirExistsInDatabase {
		log.Error().Err(err).Int("dirId", dirID).Msg("Directory not found in database")
		return false, err
	}

	absolutePath, err := u.dirService.CalcAbsolutePath(tx, dirID)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirID).Msg("Failed to calculate absolute path")
		return false, err
	}
	existOnDisk, err := u.dirService.IsExistsOnDisk(absolutePath)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirID).Msg("Failed to check directory existence on disk")
		return false, err
	}
	if !existOnDisk {
		log.Error().Err(err).Int("dirId", dirID).Msg("Directory not found on disk")
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

func (u UseCase) scanDirScanContent(tx *sqlx.Tx, dirID int) (err error) {
	err = u.actualizeAudios(tx, dirID)
	if err != nil {
		log.Error().Int("dirId", dirID).Msg("Failed to actualize audio files")
		return err
	}

	err = u.actualizeCovers(tx, dirID)
	if err != nil {
		log.Error().Int("dirId", dirID).Msg("Failed to actualize covers")
		return err
	}

	return nil
}

func (u UseCase) actualizeAudios(tx *sqlx.Tx, dirID int) (err error) {
	dirPath, err := u.dirService.CalcAbsolutePath(tx, dirID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to calculate absolute path")
		return err
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read dir's entries")
		return err
	}

	for _, entry := range entries {
		filePath := filepath.Join(dirPath, entry.Name())
		isAudio, err := u.audioService.IsAudioByPath(filePath)
		if err != nil {
			log.Error().Err(err).Msg("Failed to check entry's type")
			return err
		}
		if !isAudio {
			continue
		}

		sha256OnDisk, err := u.audioService.CalculateSHA256(filePath)
		if err != nil {
			log.Error().Int("dirId", dirID).Msg("Failed to calculate sha256")
			return err
		}

		alreadyInDatabase, err := u.audioService.IsExistsByDirAndName(tx, dirID, entry.Name())
		if err != nil {
			log.Error().Int("dirId", dirID).Str("entryName", entry.Name()).Msg("Failed to check audio existence")
			return err
		}

		if alreadyInDatabase {
			audio, err := u.audioService.GetByDirAndName(tx, dirID, entry.Name())
			if err != nil {
				log.Error().Int("dirId", dirID).Str("entryName", entry.Name()).Msg("Failed to check audio existence")
				return err
			}
			if sha256OnDisk == audio.SHA256 {
				continue
			}

			audioToUpdate, err := u.audioService.ConstructByPath(filePath)
			if err != nil {
				log.Error().Err(err).Msg("Failed to construct audio by path")
				return err
			}
			audioToUpdate.DirID = dirID
			audioToUpdate.SHA256 = sha256OnDisk

			err = u.audioService.Update(tx, audio.ID, audioToUpdate)
			if err != nil {
				log.Error().Err(err).Msg("Failed to update audio")
				return err
			}
		} else {
			audioToCreate, err := u.audioService.ConstructByPath(filePath)
			if err != nil {
				log.Error().Err(err).Msg("Failed to construct audio by path")
				return err
			}
			audioToCreate.DirID = dirID
			audioToCreate.SHA256 = sha256OnDisk

			_, err = u.audioService.Create(tx, audioToCreate)
			if err != nil {
				log.Error().Err(err).Msg("Failed to create audio")
				return err
			}
		}
	}

	audios, err := u.audioService.GetAllByDir(tx, dirID)
	if err != nil {
		log.Error().Int("dirId", dirID).Msg("Failed to get audio")
		return err
	}

	for _, audio := range audios {
		foundOnDisk := false

		for _, entry := range entries {
			filePath := filepath.Join(dirPath, entry.Name())
			isAudio, err := u.audioService.IsAudioByPath(filePath)
			if err != nil {
				return err
			}

			if isAudio && audio.Filename == entry.Name() {
				foundOnDisk = true
			}
		}

		if !foundOnDisk {
			err = u.audioService.Delete(tx, audio.ID)
			if err != nil {
				log.Error().Int("dirId", dirID).Msg("Failed to delete audio file")
				return err
			}
		}
	}

	return nil
}

func (u UseCase) actualizeCovers(tx *sqlx.Tx, dirID int) (err error) {
	dirPath, err := u.dirService.CalcAbsolutePath(tx, dirID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to calculate absolute path")
		return err
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read dir's entries")
		return err
	}

	for _, entry := range entries {
		filePath := filepath.Join(dirPath, entry.Name())
		isCover, err := u.coverService.IsCoverByPath(filePath)
		if err != nil {
			log.Error().Err(err).Msg("Failed to check entry's type")
			return err
		}
		if !isCover {
			continue
		}

		sha256OnDisk, err := u.coverService.CalculateSHA256(filePath)
		if err != nil {
			log.Error().Int("dirId", dirID).Msg("Failed to calculate sha256")
			return err
		}

		alreadyInDatabase, err := u.coverService.IsExistsByDirAndName(tx, dirID, entry.Name())
		if err != nil {
			log.Error().Int("dirId", dirID).Str("entryName", entry.Name()).Msg("Failed to check cover existence")
			return err
		}

		if alreadyInDatabase {
			cover, err := u.coverService.GetByDirAndName(tx, dirID, entry.Name())
			if err != nil {
				log.Error().Int("dirId", dirID).Str("entryName", entry.Name()).Msg("Failed to check cover existence")
				return err
			}
			if sha256OnDisk == cover.SHA256 {
				continue
			}

			coverToUpdate, err := u.coverService.ConstructByPath(filePath)
			if err != nil {
				log.Error().Err(err).Msg("Failed to construct cover by path")
				return err
			}
			coverToUpdate.DirID = dirID
			coverToUpdate.SHA256 = sha256OnDisk

			err = u.coverService.Update(tx, cover.ID, coverToUpdate)
			if err != nil {
				log.Error().Err(err).Msg("Failed to update cover")
				return err
			}
		} else {
			coverToCreate, err := u.coverService.ConstructByPath(filePath)
			if err != nil {
				log.Error().Err(err).Msg("Failed to construct cover by path")
				return err
			}
			coverToCreate.DirID = dirID
			coverToCreate.SHA256 = sha256OnDisk

			_, err = u.coverService.Create(tx, coverToCreate)
			if err != nil {
				log.Error().Err(err).Msg("Failed to create cover")
				return err
			}
		}
	}

	covers, err := u.coverService.GetAllByDir(tx, dirID)
	if err != nil {
		log.Error().Int("dirId", dirID).Msg("Failed to get cover")
		return err
	}

	for _, cover := range covers {
		foundOnDisk := false

		for _, entry := range entries {
			filePath := filepath.Join(dirPath, entry.Name())
			isCover, err := u.coverService.IsCoverByPath(filePath)
			if err != nil {
				return err
			}

			if isCover && cover.Filename == entry.Name() {
				foundOnDisk = true
			}
		}

		if !foundOnDisk {
			err = u.coverService.Delete(tx, cover.ID)
			if err != nil {
				log.Error().Int("dirId", dirID).Msg("Failed to delete cover file")
				return err
			}
		}
	}

	return nil
}
