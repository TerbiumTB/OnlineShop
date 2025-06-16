package handlers

import (
	"log"
	"orders/services"
)

type Handler struct {
	s service.OrderServicing
	l *log.Logger
}

func NewHandler(s service.OrderServicing, l *log.Logger) *Handler {
	return &Handler{s: s, l: l}
}
