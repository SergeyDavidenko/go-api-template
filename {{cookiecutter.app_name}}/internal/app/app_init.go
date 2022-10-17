package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"{{cookiecutter.app_name}}/internal/repository"
	"{{cookiecutter.app_name}}/internal/rest"
	"{{cookiecutter.app_name}}/pkg/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Server struct {
	cfg     *config.Config
	quit    chan os.Signal
	repo    *repository.DB
	api     *fiber.App
	heltz   *fiber.App
	handler *rest.Handler
}

func New(conf *config.Config, repo *repository.DB) *Server {
	initLogger()
	cfgApi := conf.GetHTTP("api")
	srv := &Server{
		api: fiber.New(fiber.Config{
			ReadTimeout:  cfgApi.ReadTimeout,
			WriteTimeout: cfgApi.WriteTimeout,
			IdleTimeout:  15 * time.Second,
		}),
		cfg:     conf,
		quit:    make(chan os.Signal, 1),
		repo:    repo,
		handler: rest.New(repo),
	}
	signal.Notify(srv.quit, syscall.SIGINT, syscall.SIGTERM)
	return srv
}

func (s *Server) Run() {
	err := s.repo.Migrations(".")
	if err != nil {
		logrus.Fatal(err)
	}
	s.api.Use(cors.New())
	s.api.Use(recover.New())
	s.setupRouter()
	if viper.GetBool("USE_HEALTH") {
		go func() {
			if errAppHealtz := s.heltz.Listen(s.cfg.GetHTTP("healtz").HostString); errAppHealtz != nil && errAppHealtz != http.ErrServerClosed {
				logrus.Fatal("listen: ", errAppHealtz)
			}
		}()
	}
	go func() {
		if err := s.api.Listen(s.cfg.GetHTTP("api").HostString); err != nil && err != http.ErrServerClosed {
			logrus.Fatal("listen: ", err)
		}
	}()
	<-s.quit
	logrus.Info("server shutdown ...")
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.shutdown()
	logrus.Info("server exiting")
}

func (s *Server) shutdown() {
	s.api.Shutdown()
}

func initLogger() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
