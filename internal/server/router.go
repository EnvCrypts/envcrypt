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
	router.Handle("/env/", http.StripPrefix("/env", EnvRouter(handler)))
	router.Handle("/service_role/", http.StripPrefix("/service_role", ServiceRoleRouter(handler)))
	router.Handle("/oidc/", http.StripPrefix("/oidc", OIDCRouter(handler)))

	return router
}

func UserRouter(handler *handlers.Handler) *http.ServeMux {

	userRouter := http.NewServeMux()

	userRouter.HandleFunc("POST /create", handler.CreateUser)
	userRouter.HandleFunc("POST /login", handler.LoginUser)
	userRouter.HandleFunc("POST /search", handler.GetUserPublicKey)

	return userRouter
}

func ProjectRouter(handler *handlers.Handler) *http.ServeMux {
	projectRouter := http.NewServeMux()

	projectRouter.HandleFunc("POST /keys", handler.GetUserProjectKeys)
	projectRouter.HandleFunc("POST /create", handler.CreateProject)
	projectRouter.HandleFunc("POST /list", handler.ListProjects)
	projectRouter.HandleFunc("POST /get", handler.GetMemberProject)
	projectRouter.HandleFunc("POST /delete", handler.DeleteProject)
	projectRouter.HandleFunc("POST /addUser", handler.AddUserToProject)
	projectRouter.HandleFunc("POST /access", handler.SetUserAccess)

	return projectRouter
}

func EnvRouter(handler *handlers.Handler) *http.ServeMux {
	envRouter := http.NewServeMux()

	envRouter.HandleFunc("POST /search", handler.GetEnv)
	envRouter.HandleFunc("POST /search/all", handler.GetEnvVersions)
	envRouter.HandleFunc("POST /create", handler.AddEnv)
	envRouter.HandleFunc("POST /update", handler.UpdateEnv)

	return envRouter
}

func ServiceRoleRouter(handler *handlers.Handler) *http.ServeMux {
	serviceRoleRouter := http.NewServeMux()

	serviceRoleRouter.HandleFunc("POST /get", handler.GetServiceRole)
	serviceRoleRouter.HandleFunc("POST /get/all", handler.ListServiceRoles)
	serviceRoleRouter.HandleFunc("POST /create", handler.CreateServiceRole)
	serviceRoleRouter.HandleFunc("POST /delete", handler.DeleteServiceRole)
	serviceRoleRouter.HandleFunc("POST /delegate", handler.DelegateAccess)
	serviceRoleRouter.HandleFunc("POST /project-keys", handler.GetProjectKeys)
	serviceRoleRouter.HandleFunc("POST /perms", handler.GetPerms)

	return serviceRoleRouter
}

func OIDCRouter(handler *handlers.Handler) *http.ServeMux {
	oidcRouter := http.NewServeMux()

	oidcRouter.HandleFunc("POST /github", handler.GitHubOIDCLogin)

	return oidcRouter
}
