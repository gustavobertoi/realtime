package server

import "github.com/gin-gonic/gin"

func (s *Server) registerRoutes(r *gin.Engine) error {
	r.Static("/chat", "./web/chat")
	r.Static("/notifications", "./web/notifications")

	r.GET("/channels/:channelId", func(ctx *gin.Context) {
		channelById(ctx, s.c)
	})

	r.POST("/channels/:channelId", func(ctx *gin.Context) {
		pushServerMessage(ctx, s.c)
	})

	return nil
}
