package app

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
	"music-files/internal/config"
	"music-files/internal/handler"
	"music-files/internal/middleware"
	"music-files/internal/service"
	"music-files/internal/service/audio_service"
	"music-files/internal/service/cover_service"
	"music-files/internal/service/dir_service"
	"music-files/internal/storage"
	"music-files/internal/storage/repo/audio_repo"
	"music-files/internal/storage/repo/cover_repo"
	"music-files/internal/storage/repo/dir_repo"
	"music-files/internal/use_case"
)

type App struct {
	handler *handler.Handler
	router  *gin.Engine
}

func New(cfg config.Config) (*App, error) {
	application := &App{}

	db, err := storage.New(cfg.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create storage")
	}

	transactor := service.NewTransactor(*db)

	audioRepo := audio_repo.New()
	coverRepo := cover_repo.New()
	dirRepo := dir_repo.New()

	audioService := audio_service.New(audioRepo)
	coverService := cover_service.New(coverRepo)
	dirService := dir_service.New(dirRepo)

	useCase := use_case.New(transactor, audioService, coverService, dirService)

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", json.Unmarshal)
	bundle.LoadMessageFile("internal/locale/en-US.json")
	bundle.LoadMessageFile("internal/locale/ru-RU.json")
	application.handler = handler.New(useCase, *bundle)

	gin.SetMode(gin.ReleaseMode)
	application.router = gin.New()
	application.router.Use(middleware.ZerologMiddleware(log.Logger))
	application.router.Use(middleware.ProduceLanguageMiddleware())
	application.RegisterRoutes()

	return application, err
}
