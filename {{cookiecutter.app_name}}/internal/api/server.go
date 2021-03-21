package api

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/{{cookiecutter.github_username}}/{{cookiecutter.app_name}}/internal/config"
	"github.com/{{cookiecutter.github_username}}/{{cookiecutter.app_name}}/api/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	{% if cookiecutter.use_postgresql == "y" %}
	"github.com/{{cookiecutter.github_username}}/{{cookiecutter.app_name}}/internal/storage"
	{% endif %}
	log "github.com/sirupsen/logrus"
)

var (
	// Quit gc shutdown
	Quit = make(chan os.Signal, 1)
)

// WebServerFiberRun and init middleware cors, recover
func WebServerFiberRun() {
	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	})
	app.Use(cors.New())
	app.Use(recover.New())
	if strings.ToUpper(config.Conf.Log.Level) == "DEBUG" {
		app.Use(logger.New(logger.Config{
			Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
			TimeFormat: "2006-01-02T15:04:05Z07:00",
			Output:     os.Stdout,
		}))
	}
	appHealtz := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	})
	appHealtz.Use(pprof.New())
	router.SetupRoutesHealtz(appHealtz)
	router.SetupRoutes(app)
	go func() {
		if errAppHealtz := appHealtz.Listen(config.Conf.API.HealthPort); errAppHealtz != nil && errAppHealtz != http.ErrServerClosed {
			log.Fatal("listen: ", errAppHealtz)
		}
	}()
	go func() {
		if err := app.Listen(config.Conf.API.Port); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: ", err)
		}
	}()
	<-Quit
	log.Info("Server shutdown ...")
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Shutdown(); err != nil {
		log.Fatal("Server shutdown: ", err)
	}
	{% if cookiecutter.use_postgresql == "y" %}
	errClose := storage.StorageDB.Close()
	if errClose != nil {
		log.Error(errClose)
	} else {
		log.Info("Close database success")
	}
	{% endif %}
	{% if cookiecutter.use_redis == "y" %}
	errCloseRedis := storage.CacheRedis.Close()
	if errCloseRedis != nil {
		log.Error(errCloseRedis)
	} else {
		log.Info("Close redis success")
	}
	{% endif %}
	log.Info("Server exiting")
}