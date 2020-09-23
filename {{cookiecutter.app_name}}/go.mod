module github.com/{{cookiecutter.github_username}}/{{cookiecutter.app_name}}

go {{cookiecutter.docker_build_image_version}}

require (
	github.com/gofiber/adaptor/v2 v2.0.0
	github.com/gofiber/fiber/v2 v2.0.2
	github.com/prometheus/client_golang v1.7.1
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/viper v1.7.1
	github.com/valyala/fasthttp v1.16.0
)