package app

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"{{cookiecutter.app_name}}/internal/repository"
	"{{cookiecutter.app_name}}/internal/rest"
	"{{cookiecutter.app_name}}/pkg/config"
)

type Server struct {
	cfg          *config.Config
	database     *repository.DB
	api          *fiber.App
	healthServer *fiber.App
	handler      *rest.Handler
}

func New(conf *config.Config, database *repository.DB) *Server {
	initLogger()
	cfgApi := conf.GetHTTP("api")
	srv := &Server{
		api: fiber.New(fiber.Config{
			ReadTimeout:  cfgApi.ReadTimeout,
			WriteTimeout: cfgApi.WriteTimeout,
			IdleTimeout:  15 * time.Second,
		}),
		healthServer: fiber.New(),
		cfg:          conf,
		database:     database,
		handler:      rest.New(database, conf),
	}
	return srv
}

func (s *Server) run() {
	if err := s.database.Migrations("."); err != nil {
		logrus.Fatalf("failed to apply migrations: %v", err)
	}
	s.api.Use(cors.New())
	s.api.Use(recover.New())
	s.setupRouter()
	if viper.GetBool("USE_HEALTH") {
		go func() {
			if err := s.healthServer.Listen(s.cfg.GetHTTP("healtz").HostString); err != nil && err != http.ErrServerClosed {
				logrus.Fatalf("Health server listen error: %v", err)
			}
		}()
	}
	go func() {
		if err := s.api.Listen(s.cfg.GetHTTP("api").HostString); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("API server listen error: %v", err)
		}
	}()
}

func (s *Server) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if s.api != nil {
		if err := s.api.ShutdownWithContext(ctx); err != nil {
			logrus.Errorf("API server shutdown error: %v", err)
		}
	}

	if s.healthServer != nil {
		if err := s.healthServer.ShutdownWithContext(ctx); err != nil {
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

func RunServer(lc fx.Lifecycle, srv *Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			srv.run()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			srv.shutdown()
			return nil
		},
	})
}
