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
	router.Handle("/projects/", http.StripPrefix("/projects", ProjectRouter(handler)))

	return router
}

func UserRouter(handler *handlers.Handler) *http.ServeMux {

	userRouter := http.NewServeMux()

	userRouter.HandleFunc("GET /all", handler.GetUsers)
	userRouter.HandleFunc("POST /create", handler.CreateUser)
	userRouter.HandleFunc("POST /login", handler.LoginUser)

	return userRouter
}

func ProjectRouter(handler *handlers.Handler) *http.ServeMux {
	projectRouter := http.NewServeMux()

	projectRouter.HandleFunc("POST /create", handler.CreateProject)

	return projectRouter
}
