package server

import "github.com/gin-gonic/gin"

func (s *Server) registerRoutes(r *gin.Engine) error {
	r.Static("/webchat", "./static")
	r.GET("/channels/:channel_id", func(ctx *gin.Context) {
		channelById(ctx, s.c)
	})
	return nil
}
