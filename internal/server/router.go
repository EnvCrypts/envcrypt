package server

import (
	"net/http"

	"github.com/vijayvenkatj/envcrypt/database"
	"github.com/vijayvenkatj/envcrypt/internal/handlers"
	"github.com/vijayvenkatj/envcrypt/internal/services"
)

func NewRouter(dbQueries *database.Queries) *http.ServeMux {
	router := http.NewServeMux()

	service := services.NewServices(dbQueries)
	handler := handlers.NewHandler(service)

	router.HandleFunc("GET /users/getall", handler.GetUsers)

	return router
}
