package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"music-files/internal/config"
	"net/http"
	"os"
	"path/filepath"
)

func Scan(c *gin.Context, cfg *config.Configuration) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	dir := filepath.Join(home, "Music")

	files, err := getFilesFromDirectory(dir)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get files from directory"})
		return
	}

	for _, file := range files {
		log.Println(file)
	}
}

func getFilesFromDirectory(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".mp3" || filepath.Ext(path) == ".flac" {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return files, nil
}
