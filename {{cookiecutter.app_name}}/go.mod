module github.com/{{cookiecutter.github_username}}/{{cookiecutter.app_name}}

go {{cookiecutter.docker_build_image_version}}

require (
	github.com/gofiber/adaptor/v2 v2.0.0
	github.com/gofiber/fiber/v2 v2.0.2
	github.com/prometheus/client_golang v1.7.1
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/viper v1.7.1
	{% if cookiecutter.use_postgresql == "y" -%}github.com/jackc/pgx v3.6.2+incompatible{%- endif %}
	{% if cookiecutter.use_postgresql == "y" -%}github.com/jmoiron/sqlx v1.2.0{%- endif %}
)