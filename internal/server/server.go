package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/vijayvenkatj/envcrypt/database"
	"github.com/vijayvenkatj/envcrypt/internal/config"
)

type Server struct {
	HttpServer *http.Server
}

func NewServer(cfg *config.Config) *Server {
	conn, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}
	dbQueries := database.New(conn)

	router := NewRouter(dbQueries)
	handler := RequestMiddleware(router)
	return &Server{
		HttpServer: &http.Server{
			Addr:    cfg.Addr,
			Handler: handler,
		},
	}
}

func (s *Server) Start() error {
	log.Printf("Server listening on %s", s.HttpServer.Addr)
	
	err := s.HttpServer.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
