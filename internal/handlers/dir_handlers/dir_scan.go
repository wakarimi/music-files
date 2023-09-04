package dir_handlers

import (
	"github.com/gin-gonic/gin"
	"music-files/internal/database/repository"
	"music-files/internal/handlers/types"
	"music-files/internal/models"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Scan(c *gin.Context) {
	dirs, err := repository.GetAllDirs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to read directories",
		})
		return
	}

	var foundMusics []models.MusicFile
	var foundCovers []models.CoverFile

	for _, dir := range dirs {
		foundMusicsOneDir, err := searchMusicsFromDirectory(dir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.Error{
				Error: "Failed to get music files from directory",
			})
			return
		}
		foundMusics = append(foundMusics, foundMusicsOneDir...)

		foundCoversOneDir, err := searchCoversFromDirectory(dir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.Error{
				Error: "Failed to get cover files from directory",
			})
			return
		}
		foundCovers = append(foundCovers, foundCoversOneDir...)
	}

	currentMusics, err := repository.GetAllMusicFiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to get full music list",
		})
		return
	}

	currentCovers, err := repository.GetAllCoverFiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to get full cover list",
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

		if isMusicFile(filepath.Ext(path)) {
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

func isMusicFile(path string) bool {
	musicExtensions := []string{".aac", ".flac", ".m4a", ".mp3", ".ogg", ".opus", ".wav", ".wma"}
	ext := strings.ToLower(filepath.Ext(path))
	for _, musicExt := range musicExtensions {
		if ext == musicExt {
			return true
		}
	}
	return false
}

func searchCoversFromDirectory(dir models.Directory) (coverFiles []models.CoverFile, err error) {
	err = filepath.Walk(dir.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if isImageFile(filepath.Ext(path)) {
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

func isImageFile(path string) bool {
	imageExtensions := []string{".jpeg", ".jpg", ".png", ".gif", ".bmp", ".tif", ".tiff", ".webp", ".ico", ".svg"}
	ext := strings.ToLower(filepath.Ext(path))
	for _, imageExt := range imageExtensions {
		if ext == imageExt {
			return true
		}
	}
	return false
}
