package handler

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"music-files/internal/use_case"
)

type useCase interface {
	AddRoot(input use_case.AddRootInput) (output use_case.AddRootOutput, err error)
	DeleteRoot(input use_case.DeleteRootInput) (output use_case.DeleteRootOutput, err error)
	GetRoots(input use_case.GetRootsInput) (output use_case.GetRootsOutput, err error)
	ScanDir(input use_case.ScanDirInput) (output use_case.ScanDirOutput, err error)
	StaticAudio(input use_case.StaticAudioInput) (output use_case.StaticAudioOutput, err error)
	StaticCover(input use_case.StaticCoverInput) (output use_case.StaticCoverOutput, err error)
	GetDirContent(input use_case.GetDirContentInput) (output use_case.GetDirContentOutput, err error)
	GetAllAudios(input use_case.GetAllAudiosInput) (output use_case.GetAllAudiosOutput, err error)
}

type Handler struct {
	useCase      useCase
	bundle       i18n.Bundle
	engLocalizer i18n.Localizer
}

func New(useCase useCase,
	bundle i18n.Bundle) *Handler {
	return &Handler{
		useCase:      useCase,
		bundle:       bundle,
		engLocalizer: *i18n.NewLocalizer(&bundle, "en_US"),
	}
}
