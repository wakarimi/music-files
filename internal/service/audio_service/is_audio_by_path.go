package audio_service

import (
	"github.com/h2non/filetype"
	"os"
)

func (s Service) IsAudioByPath(path string) (isAudio bool, err error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if fileInfo.IsDir() {
		return false, err
	}

	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	head := make([]byte, 261)
	file.Read(head)

	kind, _ := filetype.Match(head)
	if kind == filetype.Unknown {
		return false, nil
	}

	isAudio = kind.MIME.Value == "audio/mpeg" ||
		kind.MIME.Value == "audio/wav" ||
		kind.MIME.Value == "audio/flac" ||
		kind.MIME.Value == "audio/x-flac" ||
		kind.MIME.Value == "audio/aac" ||
		kind.MIME.Value == "audio/ogg" ||
		kind.MIME.Value == "audio/x-ms-wma" ||
		kind.MIME.Value == "audio/vnd.rn-realaudio" ||
		kind.MIME.Value == "audio/amr" ||
		kind.MIME.Value == "audio/mp4" ||
		kind.MIME.Value == "audio/alac" ||
		kind.MIME.Value == "audio/midi"

	return isAudio, nil
}
