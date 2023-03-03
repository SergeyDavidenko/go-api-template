module {{cookiecutter.app_name}}

go {{cookiecutter.docker_build_image_version}}

require (
	github.com/gofiber/fiber/v2 v2.42.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/viper v1.15.0
	github.com/stretchr/testify v1.8.2
	gorm.io/driver/postgres v1.3.8
	gorm.io/gorm v1.23.8
)
