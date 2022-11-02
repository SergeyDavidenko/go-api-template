package main

import (
	"{{cookiecutter.app_name}}/internal/app"
	"{{cookiecutter.app_name}}/internal/repository"
	"{{cookiecutter.app_name}}/pkg/config"
	log "github.com/sirupsen/logrus"
)

const (
	serviceName = "{{cookiecutter.app_name}}"
)

func main() {
	cfg, err := config.New("configs/", serviceName)
	if err != nil {
		log.Fatal(err)
	}
	repo, err := repository.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	srv := app.New(cfg, repo)
	srv.Run()
}
