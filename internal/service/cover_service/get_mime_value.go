package cover_service

import (
	"fmt"
	"path/filepath"
)

func (s Service) GetMimeValue(absolutePath string) (mimeValue string, err error) {
	switch filepath.Ext(absolutePath) {
	case ".jpg", ".jpeg":
		return "image/jpeg", nil
	case ".png":
		return "image/png", nil
	case ".gif":
		return "image/gif", nil
	default:
		return "", fmt.Errorf("undefind format")
	}
}
