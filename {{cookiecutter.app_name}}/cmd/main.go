package main

import (
	"{{cookiecutter.app_name}}/internal/app"
	"{{cookiecutter.app_name}}/internal/repository"
	"{{cookiecutter.app_name}}/pkg/config"

	"go.uber.org/fx"
)

const (
	serviceName = "{{cookiecutter.app_name}}"
)

func main() {
	srv := fx.New(
		fx.Provide(
			func() (*config.Config, error) {
				return config.New("configs/", serviceName)
			},
			repository.New,
			app.New,
		),
		fx.Invoke(app.RunServer),
		fx.NopLogger,
	)
	srv.Run()
}
