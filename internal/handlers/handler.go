package handlers

import (
	"context"
	"log"

	"github.com/vijayvenkatj/envcrypt/internal/services"
)

type Handler struct {
	Services *services.Services
	OIDC     OIDCVerifier
}

func NewHandler(services *services.Services) *Handler {
	verifier, err := NewGithubOIDCVerifier(
		context.Background(),
		"envcrypts/envcrypt",
	)
	if err != nil {
		log.Fatal(err)
	}

	return &Handler{
		Services: services,
		OIDC:     verifier,
	}
}
