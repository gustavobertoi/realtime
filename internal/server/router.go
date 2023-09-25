package server

import "github.com/gin-gonic/gin"

func (s *Server) registerRoutes(r *gin.Engine) error {
	r.Static("/chat", "./web/chat")

	r.GET("/channels/:channelId", func(ctx *gin.Context) {
		channelById(ctx, s.c)
	})

	return nil
}
