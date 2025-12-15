package handlers

import "github.com/vijayvenkatj/envcrypt/internal/services"

type Handler struct {
	Services *services.Services
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{
		Services: services,
	}
}
