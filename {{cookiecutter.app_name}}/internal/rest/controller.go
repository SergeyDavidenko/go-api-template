package rest

import (
	"{{cookiecutter.app_name}}/internal/repository"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	repo *repository.DB
}

func New(repository *repository.DB) *Handler {
	logrus.Debug("init handler")
	return &Handler{
		repo: repository,
	}
}
