package app

import (
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
)


func (s *Server) setupRouter() {
	s.api.Get("/version", s.handler.Version)

	configHealthcheck := healthcheck.ConfigDefault
	configHealthcheck.ReadinessEndpoint = "/readiness"
	configHealthcheck.LivenessEndpoint = "/liveness"
	s.healthServer.Use(healthcheck.New(configHealthcheck))
}
