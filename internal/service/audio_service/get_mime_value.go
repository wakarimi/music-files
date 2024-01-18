package audio_service

import (
	"fmt"
	"path/filepath"
)

func (s Service) GetMimeValue(absolutePath string) (mimeValue string, err error) {
	switch filepath.Ext(absolutePath) {
	case ".mp3":
		return "mpeg", nil
	case ".flac":
		return "flac", nil
	case ".ogg":
		return "ogg", nil
	case ".wav":
		return "wav", nil
	case ".aac":
		return "aac", nil
	case ".wma":
		return "x-ms-wma", nil
	case ".ra", ".rm":
		return "vnd.rn-realaudio", nil
	case ".amr":
		return "amr", nil
	case ".mp4":
		return "mp4", nil
	case ".m4a":
		return "alac", nil
	case ".mid", ".midi":
		return "midi", nil
	default:
		return "", fmt.Errorf("undefind format")
	}
}
