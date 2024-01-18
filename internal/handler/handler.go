package handler

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"time"
)

type useCase interface {
	AddRoot(input AddRootInput) (output AddRootOutput, err error)
	DeleteRoot(input DeleteRootInput) (output DeleteRootOutput, err error)
	GetRoots(input GetRootsInput) (output GetRootsOutput, err error)
	ScanDir(input ScanDirInput) (output ScanDirOutput, err error)
	StaticAudio(input StaticAudioInput) (output StaticAudioOutput, err error)
	StaticCover(input StaticCoverInput) (output StaticCoverOutput, err error)
	GetDirContent(input GetDirContentInput) (output GetDirContentOutput, err error)
}

type AddRootInput struct {
	Path string
}

type AddRootOutput struct {
	DirID int
	Path  string
}

type DeleteRootInput struct {
	DirID int
}

type DeleteRootOutput struct {
}

type GetRootsInput struct {
}

type GetRootsOutputDirItem struct {
	DirID       int
	Path        string
	LastScanned *time.Time
}

type GetRootsOutput struct {
	Dirs []GetRootsOutputDirItem
}

type ScanDirInput struct {
	DirID int
}

type ScanDirOutput struct{}

type StaticAudioInput struct {
	AudioID int
}

type StaticAudioOutput struct {
	AbsolutePath string
	Mime         string
}

type StaticCoverInput struct {
	CoverID int
}

type StaticCoverOutput struct {
	AbsolutePath string
	Mime         string
}

type Handler struct {
	useCase      useCase
	bundle       i18n.Bundle
	engLocalizer i18n.Localizer
}

type GetDirContentInput struct {
	DirID int
}

type GetDirContentOutputDirs struct {
	ID          int
	Name        string
	LastScanned *time.Time
}

type GetDirContentOutputAudios struct {
	ID                int
	DirID             int
	DurationMs        int64
	SHA256            string
	LastContentUpdate time.Time
}

type GetDirContentOutput struct {
	Dirs   []GetDirContentOutputDirs
	Audios []GetDirContentOutputAudios
}

func New(useCase useCase,
	bundle i18n.Bundle) *Handler {
	return &Handler{
		useCase:      useCase,
		bundle:       bundle,
		engLocalizer: *i18n.NewLocalizer(&bundle, "en_US"),
	}
}
