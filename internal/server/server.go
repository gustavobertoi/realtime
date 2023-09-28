package server

import (
	"github.com/gin-gonic/gin"
	"github.com/open-source-cloud/realtime/internal/config"
)

type Server struct {
	c *config.Config
}

func NewServer(c *config.Config) *Server {
	s := &Server{
		c: c,
	}
	return s
}

func (s *Server) Start() error {
	r := gin.New()
	r.Use(gin.Recovery())
	if err := s.registerRoutes(r); err != nil {
		return err
	}
	if err := r.Run(s.c.GetPort()); err != nil {
		return err
	}
	return nil
}
