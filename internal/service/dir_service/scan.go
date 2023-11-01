package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/wtolson/go-taglib"
	"image"
	"music-files/internal/errors"
	"music-files/internal/models"
	"music-files/internal/utils"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func (s *Service) Scan(tx *sqlx.Tx, dirId int) (err error) {
	log.Debug().Int("dirId", dirId).Msg("Scanning directory")

	existsInDatabase, err := s.DirRepo.IsExists(tx, dirId)
	if err != nil {
		log.Error().Int("dirId", dirId).Err(err).Msg("Failed to check directory existence")
		return err
	}
	if !existsInDatabase {
		log.Error().Int("dirId", dirId).Msg("Directory not found")
		return errors.NotFound{Resource: "directory in database"}
	}

	absolutePath, err := s.AbsolutePath(tx, dirId)
	if err != nil {
		log.Error().Int("dirId", dirId).Err(err).Msg("Failed to calculate absolute path to directory")
		return err
	}
	existsOnDisk, err := utils.IsDirectoryExistsOnDisk(absolutePath)
	if err != nil {
		log.Error().Int("dirId", dirId).Err(err).Msg("Failed to check directory existence on disk")
		return err
	}
	if !existsOnDisk {
		err = s.DeleteDir(tx, dirId)
		if err != nil {
			log.Error().Int("dirId", dirId).Err(err).Msg("Failed to delete directory")
			return err
		}
		return nil
	}

	err = s.actualizeSubDirs(tx, dirId)
	if err != nil {
		log.Error().Int("dirId", dirId).Err(err).Msg("Failed to actualize subdirectories")
		return err
	}

	subDirs, err := s.DirRepo.ReadSubDirs(tx, dirId)
	if err != nil {
		log.Error().Int("dirId", dirId).Err(err).Msg("Failed to read subdirectories")
		return err
	}

	for _, subDir := range subDirs {
		err = s.Scan(tx, subDir.DirId)
		if err != nil {
			log.Error().Int("subDirId", subDir.DirId).Err(err).Msg("Failed to scan subdirectory")
			return err
		}
	}
	err = s.scanContent(tx, dirId)
	if err != nil {
		log.Error().Int("dirId", dirId).Err(err).Msg("Failed to scan directory's content")
		return err
	}

	log.Debug().Int("dirId", dirId).Msg("Directory scanned successfully")
	return nil
}

func (s *Service) actualizeSubDirs(tx *sqlx.Tx, dirId int) (err error) {
	absolutePath, err := s.AbsolutePath(tx, dirId)
	if err != nil {
		log.Error().Int("dirId", dirId).Msg("failed to actualize subdirectories")
		return err
	}

	entries, err := os.ReadDir(absolutePath)
	if err != nil {
		log.Error().Str("absolutePath", absolutePath).Msg("Failed to read directory from disk")
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			alreadyInDatabase, err := s.DirRepo.IsExistsByParentAndName(tx, &dirId, entry.Name())
			if err != nil {
				log.Error().Int("dirId", dirId).Msg("Failed to check directory existence")
				return err
			}
			if !alreadyInDatabase {
				_, err = s.DirRepo.Create(tx, models.Directory{
					ParentDirId: &dirId,
					Name:        entry.Name(),
				})
				if err != nil {
					log.Error().Int("dirId", dirId).Msg("Failed to create directory")
					return err
				}
			}
		}
	}

	subDirs, err := s.DirRepo.ReadSubDirs(tx, dirId)
	if err != nil {
		log.Error().Int("dirId", dirId).Msg("Failed to read subdirectories")
		return err
	}
	for _, subDir := range subDirs {
		foundDirOnDisk := false

		for _, entry := range entries {
			if entry.IsDir() {
				if subDir.Name == entry.Name() {
					foundDirOnDisk = true
				}
			}
		}

		if !foundDirOnDisk {
			err = s.DeleteDir(tx, subDir.DirId)
			if err != nil {
				log.Error().Int("dirId", dirId).Msg("Failed to delete directory from disk")
				return err
			}
		}
	}

	return nil
}

func (s *Service) scanContent(tx *sqlx.Tx, dirId int) (err error) {
	err = s.actualizeAudioFiles(tx, dirId)
	if err != nil {
		log.Error().Int("dirId", dirId).Msg("Failed to actualize audio files")
		return err
	}

	err = s.actualizeCovers(tx, dirId)
	if err != nil {
		log.Error().Int("dirId", dirId).Msg("Failed to actualize covers")
		return err
	}

	return nil
}

func (s *Service) actualizeAudioFiles(tx *sqlx.Tx, dirId int) (err error) {
	absolutePath, err := s.AbsolutePath(tx, dirId)
	if err != nil {
		log.Error().Int("dirId", dirId).Msg("Failed to calculate absolute path to directory")
		return err
	}

	entries, err := os.ReadDir(absolutePath)
	if err != nil {
		log.Error().Int("dirId", dirId).Msg("Failed to read directory from disk")
		return err
	}

	for _, entry := range entries {
		fileAbsolutePath := filepath.Join(absolutePath, entry.Name())
		isMusicFile, err := utils.IsMusicFile(fileAbsolutePath)
		if err != nil {
			log.Error().Int("dirId", dirId).Msg("Failed to check entry's type")
			return err
		}
		if isMusicFile {
			sha256OnDisk, err := utils.CalculateSha256(fileAbsolutePath)
			if err != nil {
				log.Error().Int("dirId", dirId).Msg("Failed to calculate sha256")
				return err
			}

			alreadyInDatabase, err := s.AudioFileService.IsExistsByDirAndName(tx, dirId, entry.Name())
			if err != nil {
				log.Error().Int("dirId", dirId).Str("entryName", entry.Name()).Msg("Failed to check music file existence")
				return err
			}

			if alreadyInDatabase {
				audioFile, err := s.AudioFileService.GetByDirAndName(tx, dirId, entry.Name())
				if err != nil {
					log.Error().Int("dirId", dirId).Str("entryName", entry.Name()).Msg("Failed to check music file existence")
					return err
				}
				sha256InDatabase := audioFile.Sha256

				if sha256OnDisk == sha256InDatabase {
					continue
				}

				audioFileToUpdate, err := s.prepareAudioFileByAbsolutePath(fileAbsolutePath)
				if err != nil {
					log.Error().Int("dirId", dirId).Str("entryName", entry.Name()).Msg("Failed to prepare audio file")
					return err
				}
				audioFileToUpdate.DirId = dirId
				audioFileToUpdate.Sha256 = sha256OnDisk

				_, err = s.AudioFileService.Update(tx, audioFile.AudioFileId, audioFileToUpdate)
				if err != nil {
					log.Error().Int("dirId", dirId).Str("entryName", entry.Name()).Msg("Failed to update audio file")
					return err
				}
			} else {
				audioFileToCreate, err := s.prepareAudioFileByAbsolutePath(fileAbsolutePath)
				if err != nil {
					log.Error().Int("dirId", dirId).Str("entryName", entry.Name()).Msg("Failed to prepare audio file")
					return err
				}
				audioFileToCreate.DirId = dirId
				audioFileToCreate.Sha256 = sha256OnDisk

				_, err = s.AudioFileService.Create(tx, audioFileToCreate)
				if err != nil {
					log.Error().Int("dirId", dirId).Str("entryName", entry.Name()).Msg("Failed to create audio file")
					return err
				}
			}

		}
	}

	audioFiles, err := s.AudioFileService.GetAllByDir(tx, dirId)
	if err != nil {
		log.Error().Int("dirId", dirId).Msg("Failed to get audio files")
		return err
	}

	for _, audioFile := range audioFiles {
		foundOnDisk := false

		for _, entry := range entries {
			fileAbsolutePath := filepath.Join(absolutePath, entry.Name())
			isMusicFile, err := utils.IsMusicFile(fileAbsolutePath)
			if err != nil {
				log.Error().Str("fileAbsolutePath", fileAbsolutePath).Msg("Failed to check entry's type")
				return err
			}

			if isMusicFile {
				if audioFile.Filename == entry.Name() {
					foundOnDisk = true
				}
			}
		}

		if !foundOnDisk {
			err = s.AudioFileService.Delete(tx, audioFile.AudioFileId)
			if err != nil {
				log.Error().Int("dirId", dirId).Msg("Failed to delete audio file")
				return err
			}
		}
	}

	return nil
}

func (s *Service) prepareAudioFileByAbsolutePath(absolutePath string) (audioFile models.AudioFile, err error) {
	fileInfo, err := os.Stat(absolutePath)
	if err != nil {
		log.Error().Str("absolutePath", absolutePath).Msg("Failed to get file info")
		return models.AudioFile{}, err
	}

	fileDetails, err := taglib.Read(absolutePath)
	if err != nil {
		log.Error().Str("absolutePath", absolutePath).Msg("Failed to read file details")
		return models.AudioFile{}, err
	}

	durationMs := int64(fileDetails.Length() / time.Millisecond)

	audioFile = models.AudioFile{
		Filename:     fileInfo.Name(),
		Extension:    filepath.Ext(absolutePath),
		SizeByte:     fileInfo.Size(),
		DurationMs:   durationMs,
		BitrateKbps:  fileDetails.Bitrate(),
		SampleRateHz: fileDetails.Samplerate(),
		ChannelsN:    fileDetails.Channels(),
	}

	return audioFile, nil
}

func (s *Service) actualizeCovers(tx *sqlx.Tx, dirId int) (err error) {
	absolutePath, err := s.AbsolutePath(tx, dirId)
	if err != nil {
		log.Error().Int("dirId", dirId).Msg("Failed to calculate absolute path to directory")
		return err
	}

	entries, err := os.ReadDir(absolutePath)
	if err != nil {
		log.Error().Int("dirId", dirId).Msg("Failed to read directory on disk")
		return err
	}

	for _, entry := range entries {
		if !strings.Contains(strings.ToLower(entry.Name()), "cover") {
			continue
		}

		fileAbsolutePath := filepath.Join(absolutePath, entry.Name())
		isImageFile, err := utils.IsImageFile(fileAbsolutePath)
		if err != nil {
			log.Error().Err(err).Str("absolutePath", fileAbsolutePath).Msg("Failed to check on image")
			return err
		}
		if isImageFile {
			sha256OnDisk, err := utils.CalculateSha256(fileAbsolutePath)
			if err != nil {
				log.Error().Str("entryName", entry.Name()).Msg("Failed to calculate sha256")
				return err
			}

			alreadyInDatabase, err := s.CoverService.IsExistsByDirAndName(tx, dirId, entry.Name())
			if err != nil {
				log.Error().Str("entryName", entry.Name()).Msg("Failed to check directory existence")
				return err
			}

			if alreadyInDatabase {
				cover, err := s.CoverService.GetByDirAndName(tx, dirId, entry.Name())
				if err != nil {
					log.Error().Str("entryName", entry.Name()).Msg("Failed to check directory existence")
					return err
				}
				sha256InDatabase := cover.Sha256

				if sha256OnDisk == sha256InDatabase {
					continue
				}

				coverToUpdate, err := s.prepareCoverByAbsolutePath(absolutePath)
				if err != nil {
					log.Error().Str("entryName", entry.Name()).Msg("Failed to prepare cover")
					return err
				}
				coverToUpdate.DirId = dirId
				coverToUpdate.Sha256 = sha256OnDisk

				_, err = s.CoverService.Update(tx, cover.CoverId, coverToUpdate)
				if err != nil {
					log.Error().Str("entryName", entry.Name()).Msg("Failed to update cover")
					return err
				}
			} else {
				coverToCreate, err := s.prepareCoverByAbsolutePath(fileAbsolutePath)
				if err != nil {
					log.Error().Str("entryName", entry.Name()).Msg("Failed to prepare cover")
					return err
				}
				coverToCreate.DirId = dirId
				coverToCreate.Sha256 = sha256OnDisk

				_, err = s.CoverService.Create(tx, coverToCreate)
				if err != nil {
					log.Error().Str("entryName", entry.Name()).Msg("Failed to create cover")
					return err
				}
			}

		}
	}

	covers, err := s.CoverService.GetAllByDir(tx, dirId)
	if err != nil {
		log.Error().Int("dirId", dirId).Msg("Failed to get covers")
		return err
	}

	for _, cover := range covers {
		foundOnDisk := false

		for _, entry := range entries {
			if !strings.Contains(strings.ToLower(entry.Name()), "cover") {
				continue
			}

			fileAbsolutePath := filepath.Join(absolutePath, entry.Name())
			isImageFile, err := utils.IsImageFile(fileAbsolutePath)
			if err != nil {
				log.Error().Str("fileAbsolutePath", fileAbsolutePath).Msg("Failed to check file type")
				return err
			}

			if isImageFile {
				if cover.Filename == entry.Name() {
					foundOnDisk = true
				}
			}
		}

		if !foundOnDisk {
			err = s.CoverService.Delete(tx, cover.CoverId)
			if err != nil {
				log.Error().Err(err).Int("coverId", cover.CoverId).Msg("Failed to delete cover")
				return err
			}
		}
	}

	return nil
}

func (s *Service) prepareCoverByAbsolutePath(absolutePath string) (audioFile models.Cover, err error) {
	fileInfo, err := os.Stat(absolutePath)
	if err != nil {
		log.Error().Err(err).Str("absolutePath", absolutePath).Msg("Failed to get file info")
		return models.Cover{}, err
	}

	f, err := os.Open(absolutePath)
	if err != nil {
		log.Error().Err(err).Str("absolutePath", absolutePath).Msg("Failed to open file")
		return models.Cover{}, err
	}
	defer f.Close()

	img, _, err := image.DecodeConfig(f)
	if err != nil {
		log.Error().Err(err).Str("absolutePath", absolutePath).Msg("Failed to decode config")
		return models.Cover{}, err
	}

	audioFile = models.Cover{
		Filename:  fileInfo.Name(),
		Extension: filepath.Ext(absolutePath),
		SizeByte:  fileInfo.Size(),
		WidthPx:   img.Width,
		HeightPx:  img.Height,
	}

	return audioFile, nil
}
