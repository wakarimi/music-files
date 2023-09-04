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

	currentMusics, err := repository.GetAllMusicFilesByDirId(dir.DirId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to get music list",
		})
		return
	}

	currentCovers, err := repository.GetAllCoverFilesByDirId(dir.DirId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to get cover list",
		})
		return
	}

	for _, fm := range foundMusics {
		if !musicExistsInDB(fm, currentMusics) {
			repository.InsertMusicFile(fm)
		}
	}
	for _, cm := range currentMusics {
		if !musicExistsInList(cm, foundMusics) {
			repository.DeleteMusicFile(cm.MusicFileId)
		}
	}

	for _, fc := range foundCovers {
		if !coverExistsInDB(fc, currentCovers) {
			repository.InsertCoverFile(fc)
		}
	}
	for _, cc := range currentCovers {
		if !coverExistsInList(cc, foundCovers) {
			repository.DeleteCoverFile(cc.CoverFileId)
		}
	}
}

func musicExistsInDB(music models.MusicFile, list []models.MusicFile) bool {
	for _, m := range list {
		if m.Path == music.Path {
			return true
		}
	}
	return false
}

func coverExistsInDB(cover models.CoverFile, list []models.CoverFile) bool {
	for _, c := range list {
		if c.Path == cover.Path {
			return true
		}
	}
	return false
}

func musicExistsInList(music models.MusicFile, list []models.MusicFile) bool {
	return musicExistsInDB(music, list)
}

func coverExistsInList(cover models.CoverFile, list []models.CoverFile) bool {
	return coverExistsInDB(cover, list)
}

func searchMusicsFromDirectory(dir models.Directory) (musicFiles []models.MusicFile, err error) {
	err = filepath.Walk(dir.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if utils.IsMusicFile(filepath.Ext(path)) {
			musicFiles = append(musicFiles, models.MusicFile{
				DirId:  dir.DirId,
				Path:   path,
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

func searchCoversFromDirectory(dir models.Directory) (coverFiles []models.CoverFile, err error) {
	err = filepath.Walk(dir.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if utils.IsImageFile(filepath.Ext(path)) {
			coverFiles = append(coverFiles, models.CoverFile{
				DirId:  dir.DirId,
				Path:   path,
				Size:   info.Size(),
				Format: filepath.Ext(path),
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return coverFiles, nil
}
