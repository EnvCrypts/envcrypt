package server

import (
	"database/sql"
	"net/http"

	"github.com/vijayvenkatj/envcrypt/database"
	"github.com/vijayvenkatj/envcrypt/internal/handlers"
	"github.com/vijayvenkatj/envcrypt/internal/services"
)

func NewRouter(dbQueries *database.Queries, db *sql.DB, debug bool) *http.ServeMux {
	router := http.NewServeMux()

	auditService := services.NewAuditService(dbQueries)
	service := services.NewServices(dbQueries, auditService, db)
	handler := handlers.NewHandler(service)

	router.Handle("/users/", http.StripPrefix("/users", UserRouter(handler, debug)))
	router.Handle("/projects/", http.StripPrefix("/projects", ProjectRouter(handler, debug)))
	router.Handle("/env/", http.StripPrefix("/env", EnvRouter(handler, debug)))
	router.Handle("/service_role/", http.StripPrefix("/service_role", ServiceRoleRouter(handler, debug)))
	router.Handle("/oidc/", http.StripPrefix("/oidc", OIDCRouter(handler, debug)))

	return router
}

func UserRouter(handler *handlers.Handler, debug bool) *http.ServeMux {

	userRouter := http.NewServeMux()

	userRouter.HandleFunc("POST /create", WithErrors(debug, handler.CreateUser))
	userRouter.HandleFunc("POST /login", WithErrors(debug, handler.LoginUser))
	userRouter.HandleFunc("POST /logout", WithErrors(debug, handler.Logout))
	userRouter.HandleFunc("POST /search", WithErrors(debug, handler.GetUserPublicKey))
	userRouter.HandleFunc("POST /refresh", WithErrors(debug, handler.Refresh))
	userRouter.HandleFunc("POST /recovery/init", WithErrors(debug, handler.RecoveryInit))
	userRouter.HandleFunc("POST /recovery/complete", WithErrors(debug, handler.RecoveryComplete))

	return userRouter
}

func ProjectRouter(handler *handlers.Handler, debug bool) *http.ServeMux {
	projectRouter := http.NewServeMux()

	projectRouter.HandleFunc("POST /keys", WithErrors(debug, AuthMiddleware(handler.Services.SessionService, handler.GetUserProjectKeys)))
	projectRouter.HandleFunc("POST /create", WithErrors(debug, AuthMiddleware(handler.Services.SessionService, handler.CreateProject)))
	projectRouter.HandleFunc("POST /list", WithErrors(debug, AuthMiddleware(handler.Services.SessionService, handler.ListProjects)))
	projectRouter.HandleFunc("POST /get", WithErrors(debug, AuthMiddleware(handler.Services.SessionService, handler.GetMemberProject)))
	projectRouter.HandleFunc("POST /delete", WithErrors(debug, AuthMiddleware(handler.Services.SessionService, handler.DeleteProject)))
	projectRouter.HandleFunc("POST /addUser", WithErrors(debug, AuthMiddleware(handler.Services.SessionService, handler.AddUserToProject)))
	projectRouter.HandleFunc("POST /access", WithErrors(debug, AuthMiddleware(handler.Services.SessionService, handler.SetUserAccess)))
	projectRouter.HandleFunc("POST /rotate/init", WithErrors(debug, AuthMiddleware(handler.Services.SessionService, handler.RotateInit)))
	projectRouter.HandleFunc("POST /rotate/commit", WithErrors(debug, AuthMiddleware(handler.Services.SessionService, handler.RotateCommit)))

	projectRouter.HandleFunc("POST /snapshot/export", WithErrors(debug, AuthMiddleware(handler.Services.SessionService, handler.SnapshotExport)))
	projectRouter.HandleFunc("POST /snapshot/import", WithErrors(debug, AuthMiddleware(handler.Services.SessionService, handler.SnapshotImport)))
	projectRouter.HandleFunc("POST /audit", WithErrors(debug, AuthMiddleware(handler.Services.SessionService, handler.HandleProjectAuditLogs)))

	return projectRouter
}

func EnvRouter(handler *handlers.Handler, debug bool) *http.ServeMux {
	envRouter := http.NewServeMux()

	envRouter.HandleFunc("POST /search", WithErrors(debug, AuthMiddleware(handler.Services.SessionService, handler.GetEnv)))
	envRouter.HandleFunc("POST /search/all", WithErrors(debug, AuthMiddleware(handler.Services.SessionService, handler.GetEnvVersions)))
	envRouter.HandleFunc("POST /create", WithErrors(debug, AuthMiddleware(handler.Services.SessionService, handler.AddEnv)))
	envRouter.HandleFunc("POST /update", WithErrors(debug, AuthMiddleware(handler.Services.SessionService, handler.UpdateEnv)))

	envRouter.HandleFunc("POST /ci/search", WithErrors(debug, handler.GetCIEnv))

	return envRouter
}

func ServiceRoleRouter(handler *handlers.Handler, debug bool) *http.ServeMux {
	serviceRoleRouter := http.NewServeMux()

	serviceRoleRouter.HandleFunc("POST /get", WithErrors(debug, AuthMiddleware(handler.Services.SessionService, handler.GetServiceRole)))
	serviceRoleRouter.HandleFunc("POST /create", WithErrors(debug, AuthMiddleware(handler.Services.SessionService, handler.CreateServiceRole)))
	serviceRoleRouter.HandleFunc("POST /get/all", WithErrors(debug, AuthMiddleware(handler.Services.SessionService, handler.ListServiceRoles)))
	serviceRoleRouter.HandleFunc("POST /delete", WithErrors(debug, AuthMiddleware(handler.Services.SessionService, handler.DeleteServiceRole)))
	serviceRoleRouter.HandleFunc("POST /delegate", WithErrors(debug, AuthMiddleware(handler.Services.SessionService, handler.DelegateAccess)))
	serviceRoleRouter.HandleFunc("POST /perms", WithErrors(debug, AuthMiddleware(handler.Services.SessionService, handler.GetPerms)))

	serviceRoleRouter.HandleFunc("POST /project-keys", WithErrors(debug, handler.GetProjectKeys))

	return serviceRoleRouter
}

func OIDCRouter(handler *handlers.Handler, debug bool) *http.ServeMux {
	oidcRouter := http.NewServeMux()

	oidcRouter.HandleFunc("POST /github", WithErrors(debug, handler.GitHubOIDCLogin))

	return oidcRouter
}
