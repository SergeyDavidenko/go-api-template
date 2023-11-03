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
		heltz:   fiber.New(),
		cfg:     conf,
		quit:    make(chan os.Signal, 1),
		repo:    repo,
		handler: rest.New(repo),
	}
	signal.Notify(srv.quit, syscall.SIGINT, syscall.SIGTERM)
	return srv
}

func (s *Server) Run() {
	if err := s.repo.Migrations("."); err != nil {
		logrus.Fatalf("failed to apply migrations: %v", err)
	}
	s.api.Use(cors.New())
	s.api.Use(recover.New())
	s.setupRouter()
	if viper.GetBool("USE_HEALTH") {
		go func() {
			if err := s.heltz.Listen(s.cfg.GetHTTP("healtz").HostString); err != nil && err != http.ErrServerClosed {
				logrus.Fatalf("Health server listen error: %v", err)
			}
		}()
	}
	go func() {
		if err := s.api.Listen(s.cfg.GetHTTP("api").HostString); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("API server listen error: %v", err)
		}
	}()
	<-s.quit
	logrus.Info("shutting down the server...")
	s.shutdown()
	logrus.Info("server exited")
}

func (s *Server) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if s.api != nil {
		if err := s.api.ShutdownWithContext(ctx); err != nil {
			logrus.Errorf("API server shutdown error: %v", err)
		}
	}

	if s.heltz != nil {
		if err := s.heltz.ShutdownWithContext(ctx); err != nil {
			logrus.Errorf("Health server shutdown error: %v", err)
		}
	}
}

func initLogger() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
