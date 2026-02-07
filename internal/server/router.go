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
	router.Handle("/projects/", http.StripPrefix("/projects", AuthMiddleware(service.SessionService, ProjectRouter(handler))))
	router.Handle("/env/", http.StripPrefix("/env", AuthMiddleware(service.SessionService, EnvRouter(handler))))
	router.Handle("/service_role/", http.StripPrefix("/service_role", ServiceRoleRouter(handler)))
	router.Handle("/oidc/", http.StripPrefix("/oidc", OIDCRouter(handler)))

	return router
}

func UserRouter(handler *handlers.Handler) *http.ServeMux {

	userRouter := http.NewServeMux()

	userRouter.HandleFunc("POST /create", handler.CreateUser)
	userRouter.HandleFunc("POST /login", handler.LoginUser)
	userRouter.HandleFunc("POST /logout", handler.Logout)
	userRouter.HandleFunc("POST /search", handler.GetUserPublicKey)
	userRouter.HandleFunc("POST /refresh", handler.Refresh)

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
	envRouter.HandleFunc("POST /ci/search", handler.GetCIEnv)
	envRouter.HandleFunc("POST /search/all", handler.GetEnvVersions)
	envRouter.HandleFunc("POST /create", handler.AddEnv)
	envRouter.HandleFunc("POST /update", handler.UpdateEnv)

	return envRouter
}

func ServiceRoleRouter(handler *handlers.Handler) *http.ServeMux {
	serviceRoleRouter := http.NewServeMux()

	serviceRoleRouter.Handle("POST /get", AuthMiddleware(handler.Services.SessionService, http.HandlerFunc(handler.GetServiceRole)))
	serviceRoleRouter.Handle("POST /get/all", AuthMiddleware(handler.Services.SessionService, http.HandlerFunc(handler.ListServiceRoles)))
	serviceRoleRouter.Handle("POST /delete", AuthMiddleware(handler.Services.SessionService, http.HandlerFunc(handler.DeleteServiceRole)))
	serviceRoleRouter.Handle("POST /delegate", AuthMiddleware(handler.Services.SessionService, http.HandlerFunc(handler.DelegateAccess)))
	serviceRoleRouter.Handle("POST /perms", AuthMiddleware(handler.Services.SessionService, http.HandlerFunc(handler.GetPerms)))

	serviceRoleRouter.HandleFunc("POST /project-keys", handler.GetProjectKeys)

	return serviceRoleRouter
}

func OIDCRouter(handler *handlers.Handler) *http.ServeMux {
	oidcRouter := http.NewServeMux()

	oidcRouter.HandleFunc("POST /github", handler.GitHubOIDCLogin)

	return oidcRouter
}
