package directory

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"music-files/internal/handlers/types"
	"music-files/internal/models"
	"music-files/internal/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (h *Handler) Scan(c *gin.Context) {
	log.Debug().Msg("Creating new directory")

	dirIdStr := c.Param("dirId")
	dirId, err := strconv.Atoi(dirIdStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid dirId format")
		c.JSON(http.StatusBadRequest, types.Error{
			Error: "Invalid dirId format",
		})
		return
	}
	log.Debug().Int("dirId", dirId).Msg("Url parameter read successfully")

	dir, err := h.DirRepo.Read(dirId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get directory")
		c.JSON(http.StatusBadRequest, types.Error{
			Error: "Failed to get directory",
		})
		return
	}
	log.Debug().Str("path", dir.Path).Msg("Directory read successfully")

	h.dirScan(c, dir)

	log.Debug().Int("dirId", dirId).Msg("Directory scanned successfully")
	c.Status(http.StatusOK)
}

func (h *Handler) dirScan(c *gin.Context, dir models.Directory) {
	foundTracks, err := h.searchTracksFromDirectory(dir)
	if err != nil {
		log.Error().Err(err).Int("dirId", dir.DirId).Msg("Failed to get tracks from directory")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to get tracks from directory",
		})
		return
	}
	log.Debug()

	foundCovers, err := h.searchCoversFromDirectory(dir)
	if err != nil {
		log.Error().Err(err).Int("dirId", dir.DirId).Msg("Failed to get covers from directory")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to get covers from directory",
		})
		return
	}

	databaseTracks, err := h.TrackRepo.ReadAllByDirId(dir.DirId)
	if err != nil {
		log.Error().Err(err).Int("dirId", dir.DirId).Msg("Failed to get tracks from database")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to get tracks from database",
		})
		return
	}

	databaseCovers, err := h.CoverRepo.ReadAllByDirId(dir.DirId)
	if err != nil {
		log.Error().Err(err).Int("dirId", dir.DirId).Msg("Failed to get covers from database")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to get covers from database",
		})
		return
	}

	for _, databaseTrack := range databaseTracks {
		if !h.isTrackInList(databaseTrack, foundTracks) {
			err = h.TrackRepo.Delete(databaseTrack.TrackId)
			if err != nil {
				log.Error().Err(err).Int("databaseTrackId", databaseTrack.TrackId).Msg("Failed to delete track")
			} else {
				log.Info().Int("trackId", databaseTrack.TrackId).Str("databaseTrackFilename", databaseTrack.Filename).Msg("Undiscovered track deleted")
			}
		}
	}

	for _, databaseCover := range databaseCovers {
		if !h.isCoverInList(databaseCover, databaseCovers) {
			err = h.CoverRepo.Delete(databaseCover.CoverId)
			if err != nil {
				log.Error().Err(err).Int("databaseCoverId", databaseCover.CoverId).Msg("Failed to delete cover")
			} else {
				log.Info().Int("coverId", databaseCover.CoverId).Str("databaseCoverFilename", databaseCover.Filename).Msg("Undiscovered cover deleted")
			}
		}
	}
	for _, foundCover := range foundCovers {
		if !h.isCoverInList(foundCover, databaseCovers) {
			coverId, err := h.CoverRepo.Create(foundCover)
			if err != nil {
				log.Error().Err(err).Str("foundCoverRelativePath", foundCover.RelativePath).Msg("Failed to create cover")
			} else {
				log.Info().Int("coverId", coverId).Str("filename", foundCover.Filename).Msg("New cover added to database")
			}
		}
	}

	for _, foundTrack := range foundTracks {
		if !h.isTrackInList(foundTrack, databaseTracks) {
			cover, err := h.CoverRepo.ReadByDirIdAndRelativePath(dir.DirId, foundTrack.RelativePath)
			if err != nil {
				log.Error().Err(err).Int("dirId", dir.DirId).Str("relativePath", foundTrack.RelativePath).Msg("Failed to find relative cover")
			} else {
				foundTrack.CoverId = &cover.CoverId
			}

			trackId, err := h.TrackRepo.Create(foundTrack)
			if err != nil {
				log.Error().Err(err).Str("filename", foundTrack.Filename).Msg("Failed to create track")
			} else {
				log.Info().Int("trackId", trackId).Str("filename", foundTrack.Filename).Msg("New track added to database")
			}
		}
	}

	err = h.DirRepo.UpdateLastScanned(dir.DirId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update date of last scanning")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to update date of last scanning",
		})
		return
	}
}

func (h *Handler) searchTracksFromDirectory(dir models.Directory) (tracks []models.Track, err error) {
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

func (h *Handler) searchCoversFromDirectory(dir models.Directory) (covers []models.Cover, err error) {
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

func (h *Handler) isTrackInList(track models.Track, trackList []models.Track) bool {
	for _, trackListItem := range trackList {
		if trackListItem.DirId == track.DirId && trackListItem.RelativePath == track.RelativePath {
			return true
		}
	}
	return false
}

func (h *Handler) isCoverInList(cover models.Cover, coverList []models.Cover) bool {
	for _, coverListItem := range coverList {
		if coverListItem.DirId == cover.DirId && coverListItem.RelativePath == cover.RelativePath {
			return true
		}
	}
	return false
}
