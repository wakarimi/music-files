package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
	"music-files/internal/utils"
	"os"
	"path/filepath"
	"strings"
)

func (s *Service) Scan(tx *sqlx.Tx, dirId int) (err error) {
	dir, err := s.DirRepo.ReadTx(tx, dirId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get directory")
		return err
	}
	log.Debug().Str("path", dir.Path).Msg("Directory read successfully")

	s.dirScan(tx, dir)
	return nil
}

func (s *Service) dirScan(tx *sqlx.Tx, dir models.Directory) {
	foundTracks, err := s.searchTracksFromDirectory(dir)
	if err != nil {
		log.Error().Err(err).Int("dirId", dir.DirId).Msg("Failed to get tracks from directory")
		return
	}
	log.Debug().Int("foundTracksCount", len(foundTracks)).Msg("Tracks read from directory")

	foundCovers, err := s.searchCoversFromDirectory(dir)
	if err != nil {
		log.Error().Err(err).Int("dirId", dir.DirId).Msg("Failed to get covers from directory")
		return
	}
	log.Debug().Int("foundCoversCount", len(foundCovers)).Msg("Covers read from directory")

	databaseTracks, err := s.TrackRepo.ReadAllByDirIdTx(tx, dir.DirId)
	if err != nil {
		log.Error().Err(err).Int("dirId", dir.DirId).Msg("Failed to get tracks from database")
		return
	}
	log.Debug().Int("foundTracksCount", len(databaseTracks)).Msg("Tracks read from database")

	databaseCovers, err := s.CoverRepo.ReadAllByDirIdTx(tx, dir.DirId)
	if err != nil {
		log.Error().Err(err).Int("dirId", dir.DirId).Msg("Failed to get covers from database")
		return
	}
	log.Debug().Int("foundCoversCount", len(databaseCovers)).Msg("Covers read from database")

	deletedCoverCount := 0
	for _, databaseCover := range databaseCovers {
		analogFound := false
		for i := range foundCovers {
			if (databaseCover.Hash == foundCovers[i].Hash) ||
				((databaseCover.DirId == foundCovers[i].DirId) && (databaseCover.RelativePath == foundCovers[i].RelativePath) && (databaseCover.Filename == foundCovers[i].Filename)) {
				analogFound = true
				foundCovers[i].CoverId = databaseCover.CoverId
				break
			}
		}

		if !analogFound {
			err = s.TrackRepo.ResetCoverTx(tx, databaseCover.CoverId)
			if err != nil {
				log.Error().Msg("Failed to reset cover for tracks")
			} else {
				err = s.CoverRepo.DeleteTx(tx, databaseCover.CoverId)
				if err != nil {
					log.Error().Err(err).Int("databaseCoverId", databaseCover.CoverId).Msg("Failed to delete cover")
				} else {
					log.Info().Int("trackId", databaseCover.CoverId).Str("databaseCoverFilename", databaseCover.Filename).Msg("Undiscovered cover deleted")
					deletedCoverCount++
				}
			}
		}
	}
	log.Debug().Int("deletedCoversCount", deletedCoverCount).Msg("Undiscovered covers deleted")

	deletedTracksCount := 0
	for _, databaseTrack := range databaseTracks {
		analogFound := false
		for i := range foundTracks {
			if (databaseTrack.Hash == foundTracks[i].Hash) ||
				((databaseTrack.DirId == foundTracks[i].DirId) && (databaseTrack.RelativePath == foundTracks[i].RelativePath) && (databaseTrack.Filename == foundTracks[i].Filename)) {
				analogFound = true
				foundTracks[i].TrackId = databaseTrack.TrackId
				break
			}
		}

		if !analogFound {
			err = s.TrackRepo.DeleteTx(tx, databaseTrack.TrackId)
			if err != nil {
				log.Error().Err(err).Int("databaseTrackId", databaseTrack.TrackId).Msg("Failed to delete track")
			} else {
				log.Info().Int("trackId", databaseTrack.TrackId).Str("databaseTrackFilename", databaseTrack.Filename).Msg("Undiscovered track deleted")
				deletedTracksCount++
			}
		}
	}
	log.Debug().Int("deletedTracksCount", deletedTracksCount).Msg("Undiscovered tracks deleted")

	newCoversCount := 0
	modifiedCoversCount := 0
	for _, foundCover := range foundCovers {
		if foundCover.CoverId == 0 {
			coverId, err := s.CoverRepo.CreateTx(tx, foundCover)
			if err != nil {
				log.Error().Err(err).Str("relativePath", foundCover.RelativePath).Msg("Failed to create cover")
			} else {
				log.Info().Int("coverId", coverId).Str("relativePath", foundCover.RelativePath).Msg("New cover added to database")
				newCoversCount++
			}
		} else {
			err = s.CoverRepo.UpdateTx(tx, foundCover.CoverId, foundCover)
			if err != nil {
				log.Error().Err(err).Str("relativePath", foundCover.RelativePath).Msg("Failed to update cover")
			} else {
				log.Info().Int("coverId", foundCover.CoverId).Str("relativePath", foundCover.RelativePath).Msg("Cover updated in database")
				modifiedCoversCount++
			}
		}
	}
	log.Debug().Int("newCoversCount", newCoversCount).Msg("New covers added")
	log.Debug().Int("modifiedCoversCount", modifiedCoversCount).Msg("Covers modified")

	newTracksCount := 0
	modifiedTracksCount := 0
	for _, foundTrack := range foundTracks {
		cover, err := s.CoverRepo.ReadByDirIdAndRelativePathTx(tx, dir.DirId, foundTrack.RelativePath)
		if err == nil {
			foundTrack.CoverId = &cover.CoverId
		}

		if foundTrack.TrackId == 0 {
			trackId, err := s.TrackRepo.CreateTx(tx, foundTrack)
			if err != nil {
				log.Error().Err(err).Str("filename", foundTrack.Filename).Msg("Failed to create track")
			} else {
				log.Info().Int("trackId", trackId).Str("filename", foundTrack.Filename).Msg("New track added to database")
				newTracksCount++
			}
		} else {
			err = s.TrackRepo.UpdateTx(tx, foundTrack.TrackId, foundTrack)
			if err != nil {
				log.Error().Err(err).Str("filename", foundTrack.Filename).Str("relativePath", foundTrack.RelativePath).Msg("Failed to update track")
			} else {
				log.Info().Int("trackId", foundTrack.TrackId).Str("filename", foundTrack.Filename).Str("relativePath", foundTrack.RelativePath).Msg("Track updated in database")
				modifiedTracksCount++
			}
		}
	}
	log.Debug().Int("newTracksCount", newTracksCount).Msg("New tracks added")
	log.Debug().Int("modifiedTracksCount", modifiedTracksCount).Msg("Tracks Modified")

	err = s.DirRepo.UpdateLastScannedTx(tx, dir.DirId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update date of last scanning")
		return
	}
}

func (s *Service) searchTracksFromDirectory(dir models.Directory) (tracks []models.Track, err error) {
	err = filepath.Walk(dir.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if utils.IsMusicFile(filepath.Ext(path)) {
			relativeDir := filepath.Dir(strings.TrimPrefix(path, dir.Path))
			fileHash, err := utils.CalculateFileHash(path)
			if err != nil {
				log.Error().Err(err).Str("filepath", path).Msg("Failed to calculate file hash")
				return err
			}

			tracks = append(tracks, models.Track{
				DirId:        dir.DirId,
				RelativePath: relativeDir,
				Filename:     info.Name(),
				Extension:    filepath.Ext(path),
				Size:         info.Size(),
				Hash:         fileHash,
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return tracks, nil
}

func (s *Service) searchCoversFromDirectory(dir models.Directory) (covers []models.Cover, err error) {
	err = filepath.Walk(dir.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if utils.IsImageFile(filepath.Ext(path)) {
			relativeDir := filepath.Dir(strings.TrimPrefix(path, dir.Path))
			fileHash, err := utils.CalculateFileHash(path)
			if err != nil {
				log.Error().Err(err).Str("filepath", path).Msg("Failed to calculate file hash")
				return err
			}

			covers = append(covers, models.Cover{
				DirId:        dir.DirId,
				RelativePath: relativeDir,
				Filename:     info.Name(),
				Extension:    filepath.Ext(path),
				Size:         info.Size(),
				Hash:         fileHash,
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return covers, nil
}
