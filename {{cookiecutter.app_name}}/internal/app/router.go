package app

func (s *Server) setupRouter() {
	s.api.Get("/version", s.handler.Version)
}
