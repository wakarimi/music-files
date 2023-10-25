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
		return err
	}
	if !existsInDatabase {
		return errors.NotFound{Resource: "directory in database"}
	}

	absolutePath, err := s.AbsolutePath(tx, dirId)
	if err != nil {
		return err
	}
	existsOnDisk, err := utils.IsDirectoryExistsOnDisk(absolutePath)
	if err != nil {
		return err
	}
	if !existsOnDisk {
		err = s.DeleteDir(tx, dirId)
		if err != nil {
			return err
		}
		return nil
	}

	err = s.actualizeSubDirs(tx, dirId)
	if err != nil {
		return err
	}

	subDirs, err := s.DirRepo.ReadSubDirs(tx, dirId)
	if err != nil {
		return err
	}

	for _, subDir := range subDirs {
		err = s.Scan(tx, subDir.DirId)
		if err != nil {
			return err
		}
	}
	err = s.scanContent(tx, dirId)
	if err != nil {
		return err
	}

	log.Debug().Int("dirId", dirId).Msg("Directory scanned successfully")
	return nil
}

func (s *Service) actualizeSubDirs(tx *sqlx.Tx, dirId int) (err error) {
	absolutePath, err := s.AbsolutePath(tx, dirId)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(absolutePath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			alreadyInDatabase, err := s.DirRepo.IsExistsByParentAndName(tx, &dirId, entry.Name())
			if err != nil {
				return err
			}
			if !alreadyInDatabase {
				_, err = s.DirRepo.Create(tx, models.Directory{
					ParentDirId: &dirId,
					Name:        entry.Name(),
				})
				if err != nil {
					return err
				}
			}
		}
	}

	subDirs, err := s.DirRepo.ReadSubDirs(tx, dirId)
	if err != nil {
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
				return err
			}
		}
	}

	return nil
}

func (s *Service) scanContent(tx *sqlx.Tx, dirId int) (err error) {
	err = s.actualizeSongs(tx, dirId)
	if err != nil {
		return err
	}

	err = s.actualizeCovers(tx, dirId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) actualizeSongs(tx *sqlx.Tx, dirId int) (err error) {
	absolutePath, err := s.AbsolutePath(tx, dirId)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(absolutePath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fileAbsolutePath := filepath.Join(absolutePath, entry.Name())
		isMusicFile, err := utils.IsMusicFile(fileAbsolutePath)
		if err != nil {
			return err
		}
		if isMusicFile {
			sha256OnDisk, err := utils.CalculateSha256(fileAbsolutePath)
			if err != nil {
				return err
			}

			alreadyInDatabase, err := s.SongService.IsExistsByDirAndName(tx, dirId, entry.Name())
			if err != nil {
				return err
			}

			if alreadyInDatabase {
				song, err := s.SongService.GetByDirAndName(tx, dirId, entry.Name())
				if err != nil {

				}
				sha256InDatabase := song.Sha256

				if sha256OnDisk == sha256InDatabase {
					continue
				}

				songToUpdate, err := s.prepareSongByAbsolutePath(absolutePath)
				if err != nil {
					return err
				}
				songToUpdate.DirId = dirId
				songToUpdate.Sha256 = sha256OnDisk

				_, err = s.SongService.Update(tx, song.SongId, songToUpdate)
				if err != nil {
					return err
				}
			} else {
				songToCreate, err := s.prepareSongByAbsolutePath(fileAbsolutePath)
				if err != nil {
					return err
				}
				songToCreate.DirId = dirId
				songToCreate.Sha256 = sha256OnDisk

				_, err = s.SongService.Create(tx, songToCreate)
				if err != nil {
					return err
				}
			}

		}
	}

	songs, err := s.SongService.GetAllByDir(tx, dirId)
	if err != nil {
		return err
	}

	for _, song := range songs {
		foundOnDisk := false

		for _, entry := range entries {
			fileAbsolutePath := filepath.Join(absolutePath, entry.Name())
			isMusicFile, err := utils.IsMusicFile(fileAbsolutePath)
			if err != nil {
				return err
			}

			if isMusicFile {
				if song.Filename == entry.Name() {
					foundOnDisk = true
				}
			}
		}

		if !foundOnDisk {
			err = s.SongService.Delete(tx, song.SongId)
		}
	}

	subDirs, err := s.DirRepo.ReadSubDirs(tx, dirId)
	if err != nil {
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
				return err
			}
		}
	}

	return nil
}

func (s *Service) prepareSongByAbsolutePath(absolutePath string) (song models.Song, err error) {
	fileInfo, err := os.Stat(absolutePath)
	if err != nil {
		return models.Song{}, err
	}

	fileDetails, err := taglib.Read(absolutePath)
	if err != nil {
		return models.Song{}, err
	}

	durationMs := int64(fileDetails.Length() / time.Millisecond)

	song = models.Song{
		Filename:     fileInfo.Name(),
		Extension:    filepath.Ext(absolutePath),
		SizeByte:     fileInfo.Size(),
		DurationMs:   durationMs,
		BitrateKbps:  fileDetails.Bitrate(),
		SampleRateHz: fileDetails.Samplerate(),
		ChannelsN:    fileDetails.Channels(),
	}

	return song, nil
}

func (s *Service) actualizeCovers(tx *sqlx.Tx, dirId int) (err error) {
	absolutePath, err := s.AbsolutePath(tx, dirId)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(absolutePath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !strings.Contains(entry.Name(), "cover") {
			continue
		}

		fileAbsolutePath := filepath.Join(absolutePath, entry.Name())
		isImageFile, err := utils.IsImageFile(fileAbsolutePath)
		if err != nil {
			log.Warn().Err(err).Str("absolutePath", fileAbsolutePath).Msg("Failed to check on image")
			return err
		}
		if isImageFile {
			sha256OnDisk, err := utils.CalculateSha256(fileAbsolutePath)
			if err != nil {
				return err
			}

			alreadyInDatabase, err := s.CoverService.IsExistsByDirAndName(tx, dirId, entry.Name())
			if err != nil {
				return err
			}

			if alreadyInDatabase {
				cover, err := s.CoverService.GetByDirAndName(tx, dirId, entry.Name())
				if err != nil {

				}
				sha256InDatabase := cover.Sha256

				if sha256OnDisk == sha256InDatabase {
					continue
				}

				coverToUpdate, err := s.prepareCoverByAbsolutePath(absolutePath)
				if err != nil {
					return err
				}
				coverToUpdate.DirId = dirId
				coverToUpdate.Sha256 = sha256OnDisk

				_, err = s.CoverService.Update(tx, cover.CoverId, coverToUpdate)
				if err != nil {
					return err
				}
			} else {
				coverToCreate, err := s.prepareCoverByAbsolutePath(fileAbsolutePath)
				if err != nil {
					return err
				}
				coverToCreate.DirId = dirId
				coverToCreate.Sha256 = sha256OnDisk

				_, err = s.CoverService.Create(tx, coverToCreate)
				if err != nil {
					return err
				}
			}

		}
	}

	covers, err := s.CoverService.GetAllByDir(tx, dirId)
	if err != nil {
		return err
	}

	for _, cover := range covers {
		foundOnDisk := false

		for _, entry := range entries {
			if !strings.Contains(entry.Name(), "cover") {
				continue
			}

			fileAbsolutePath := filepath.Join(absolutePath, entry.Name())
			isImageFile, err := utils.IsImageFile(fileAbsolutePath)
			if err != nil {
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
		}
	}

	subDirs, err := s.DirRepo.ReadSubDirs(tx, dirId)
	if err != nil {
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
				return err
			}
		}
	}

	return nil
}

func (s *Service) prepareCoverByAbsolutePath(absolutePath string) (song models.Cover, err error) {
	fileInfo, err := os.Stat(absolutePath)
	if err != nil {
		return models.Cover{}, err
	}

	f, err := os.Open(absolutePath)
	if err != nil {
		return models.Cover{}, err
	}
	defer f.Close()

	img, _, err := image.DecodeConfig(f)
	if err != nil {
		return models.Cover{}, err
	}

	song = models.Cover{
		Filename:  fileInfo.Name(),
		Extension: filepath.Ext(absolutePath),
		SizeByte:  fileInfo.Size(),
		WidthPx:   img.Width,
		HeightPx:  img.Height,
	}

	return song, nil
}
