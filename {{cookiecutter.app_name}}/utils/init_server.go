package utils

import (
	"os/signal"
	"syscall"

	"github.com/{{cookiecutter.github_username}}/{{cookiecutter.app_name}}/api"
	"github.com/{{cookiecutter.github_username}}/{{cookiecutter.app_name}}/config"
	log "github.com/sirupsen/logrus"
)

// ServerInit init vars, config and log level
func ServerInit(configPath string) {
	var err error
	signal.Notify(api.Quit, syscall.SIGINT, syscall.SIGTERM)
	config.Conf, err = config.LoadConf(configPath)
	if err != nil {
		log.Error(err)
	}
	level, err := log.ParseLevel(config.Conf.Log.Level)
	if err != nil {
		log.Error("Cannot parse log level")
		log.SetLevel(log.InfoLevel)
	}
	log.Debug("Set log level ", level)
	log.SetLevel(level)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}