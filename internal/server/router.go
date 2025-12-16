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

	router.Handle("/users/", http.StripPrefix("/users", UserRouter(handler)))

	return router
}

func UserRouter(handler *handlers.Handler) *http.ServeMux {

	userRouter := http.NewServeMux()

	userRouter.HandleFunc("GET /all", handler.GetUsers)
	userRouter.HandleFunc("POST /create", handler.CreateUser)

	return userRouter
}
