package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/open-source-cloud/realtime/internal/config"
)

type Server struct {
	c *config.Config
}

func NewServer(config *config.Config) *Server {
	s := &Server{
		c: config,
	}
	return s
}

func (s *Server) Start() error {
	r := gin.New()
	r.Use(gin.Recovery())

	err := s.registerRoutes(r)
	if err != nil {
		return err
	}

	port := fmt.Sprintf(":%d", s.c.Port)

	err = r.Run(port)
	if err != nil {
		return err
	}

	return nil
}
