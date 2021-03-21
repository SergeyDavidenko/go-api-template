package main

import (
	"github.com/{{cookiecutter.github_username}}/{{cookiecutter.app_name}}/internal/api"
	"github.com/{{cookiecutter.github_username}}/{{cookiecutter.app_name}}/internal/app"
	"github.com/{{cookiecutter.github_username}}/{{cookiecutter.app_name}}/internal/config"
	log "github.com/sirupsen/logrus"
)


func main() {
	log.Info("Start init server")
	app.ServerInit("")
	log.Info("End init server")
	log.Infof("Start %s on port %s", "{{cookiecutter.app_name}}", config.Conf.API.Port)
	api.WebServerFiberRun()
}