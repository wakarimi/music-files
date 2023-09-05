package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"music-files/internal/database/repository"
	"music-files/internal/handlers/types"
	"music-files/internal/models"
	"music-files/internal/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func DirScan(c *gin.Context) {
	dirIdStr := c.Param("dirId")

	dirId, err := strconv.Atoi(dirIdStr)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, types.Error{
			Error: "Invalid dirId format",
		})
		return
	}

	dir, err := repository.GetDirById(dirId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, types.Error{
			Error: "Failed to get directory",
		})
		return
	}

	dirScanOne(c, dir)
}

func dirScanOne(c *gin.Context, dir models.Directory) {
	foundMusics, err := searchMusicsFromDirectory(dir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to get music files from directory",
		})
		return
	}

	foundCovers, err := searchCoversFromDirectory(dir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to get cover files from directory",
		})
		return
	}

	currentMusics, err := repository.GetTracksByDirId(dir.DirId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to get music list",
		})
		return
	}

	currentCovers, err := repository.GetAllCoversByDirId(dir.DirId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to get cover list",
		})
		return
	}

	for _, fc := range foundCovers {
		if !coverExistsInDB(fc, currentCovers) {
			_, err = repository.InsertCover(fc)
		}
	}
	for _, cc := range currentCovers {
		if !coverExistsInList(cc, foundCovers) {
			err = repository.DeleteCover(cc.CoverId)
		}
	}

	for _, fm := range foundMusics {
		if !musicExistsInDB(fm, currentMusics) {
			cover, err := repository.GetCoverByPath(fm.Path)
			if err == nil {
				fm.CoverId = &cover.CoverId
			}
			_, err = repository.InsertTrack(fm)
		}
	}
	for _, cm := range currentMusics {
		if !musicExistsInList(cm, foundMusics) {
			err = repository.DeleteTrackById(cm.TrackId)
		}
	}

	err = repository.UpdateLastScanned(dir.DirId)
	if err != nil {
		log.Println("Failed to update last scanned date:", err)
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to update last scanned date",
		})
		return
	}
}

func musicExistsInDB(music models.Track, list []models.Track) bool {
	for _, m := range list {
		if m.DirId == music.DirId && m.Path == music.Path {
			return true
		}
	}
	return false
}

func coverExistsInDB(cover models.Cover, list []models.Cover) bool {
	for _, c := range list {
		if c.DirId == cover.DirId && c.Path == cover.Path {
			return true
		}
	}
	return false
}

func musicExistsInList(music models.Track, list []models.Track) bool {
	return musicExistsInDB(music, list)
}

func coverExistsInList(cover models.Cover, list []models.Cover) bool {
	return coverExistsInDB(cover, list)
}

func searchMusicsFromDirectory(dir models.Directory) (musicFiles []models.Track, err error) {
	err = filepath.Walk(dir.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relativeDir := filepath.Dir(strings.TrimPrefix(path, dir.Path))
		if utils.IsMusicFile(filepath.Ext(path)) {
			musicFiles = append(musicFiles, models.Track{
				DirId:  dir.DirId,
				Path:   relativeDir,
				Name:   info.Name(),
				Size:   info.Size(),
				Format: filepath.Ext(path),
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return musicFiles, nil
}

func searchCoversFromDirectory(dir models.Directory) (covers []models.Cover, err error) {
	err = filepath.Walk(dir.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relativeDir := filepath.Dir(strings.TrimPrefix(path, dir.Path))
		if utils.IsImageFile(filepath.Ext(path)) {
			covers = append(covers, models.Cover{
				DirId:  dir.DirId,
				Path:   relativeDir,
				Name:   info.Name(),
				Size:   info.Size(),
				Format: filepath.Ext(path),
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return covers, nil
}
