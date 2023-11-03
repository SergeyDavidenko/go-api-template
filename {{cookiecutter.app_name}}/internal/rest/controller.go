package rest

import (
	"{{cookiecutter.app_name}}/internal/repository"
	
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	repo *repository.DB
	conf *config.Config
}

func New(repository *repository.DB, conf *config.Config) *Handler {
	log.Debug("init handler")
	return &Handler{
		repo: repository,
		conf: conf,
	}
}
