package app

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
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
	logger       *logrus.Logger
}

// customErrorHandler provides custom error handling
func customErrorHandler(log *logrus.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		log.WithError(err).WithFields(logrus.Fields{
			"path":   c.Path(),
			"method": c.Method(),
			"status": code,
		}).Error("Request error")

		return c.Status(code).JSON(fiber.Map{
			"error":  err.Error(),
			"path":   c.Path(),
			"method": c.Method(),
			"code":   code,
		})
	}
}

func New(conf *config.Config, database *repository.DB) *Server {
	log := initLogger()
	cfgApi := conf.GetHTTP("api")
	srv := &Server{
		api: fiber.New(fiber.Config{
			ReadTimeout:           cfgApi.ReadTimeout,
			WriteTimeout:          cfgApi.WriteTimeout,
			IdleTimeout:           15 * time.Second,
			ErrorHandler:          customErrorHandler(log),
			DisableStartupMessage: true,
		}),
		healthServer: fiber.New(),
		cfg:          conf,
		database:     database,
		handler:      rest.New(database, conf),
		logger:       log,
	}
	srv.setupMiddleware()
	return srv
}

// setupMiddleware configures all middleware for the API server
func (s *Server) setupMiddleware() {
	// Request ID for tracing
	s.api.Use(requestid.New())

	// CORS
	s.api.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Recovery middleware
	s.api.Use(recover.New())

	// Structured logging
	s.api.Use(logger.New(logger.Config{
		Format:     "${time} ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		Output:     s.logger.Writer(),
	}))
}

func (s *Server) run() {
	// Run database migrations
	if err := s.database.Migrations("migrations"); err != nil {
		logrus.Fatalf("failed to apply migrations: %v", err)
	}
	// Setup routes
	s.setupRouter()

	// Start health server if enabled
	if viper.GetBool("USE_HEALTH") {
		go func() {
			if err := s.healthServer.Listen(s.cfg.GetHTTP("healtz").HostString); err != nil && err != http.ErrServerClosed {
				logrus.Fatalf("Health server listen error: %v", err)
			}
		}()
	}

	// Start API server
	go func() {
		logrus.Info("Starting API server on ", s.cfg.GetHTTP("api").HostString)
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
