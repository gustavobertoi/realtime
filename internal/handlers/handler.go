package handlers

import (
	"github.com/gustavobertoi/realtime/internal/config"
	"github.com/gustavobertoi/realtime/pkg/logs"
)

type Handler struct {
	Conf   *config.Config
	Logger *logs.Logger
}

func NewHandler(conf *config.Config) *Handler {
	return &Handler{
		Conf:   conf,
		Logger: logs.NewLogger("Realtime API"),
	}
}
