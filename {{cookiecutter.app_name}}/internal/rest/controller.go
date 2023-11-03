package rest

import (
	"{{cookiecutter.app_name}}/internal/repository"
	
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	repo *repository.DB
}

func New(repository *repository.DB) *Handler {
	log.Debug("init handler")
	return &Handler{
		repo: repository,
	}
}
