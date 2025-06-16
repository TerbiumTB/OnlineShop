package handler

import (
	"log"
	"payments/services"
)

type Handler struct {
	//payments service.PaymentServicing
	as service.AccountServicing
	lg *log.Logger
}

func NewHandler(as service.AccountServicing, lg *log.Logger) *Handler {
	return &Handler{as: as, lg: lg}
}
