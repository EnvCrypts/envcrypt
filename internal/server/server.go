package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/vijayvenkatj/envcrypt/database"
	"github.com/vijayvenkatj/envcrypt/internal/config"
	dbdriver "github.com/vijayvenkatj/envcrypt/internal/db/driver"
)

type Server struct {
	HttpServer *http.Server
}

func NewServer(cfg *config.Config) *Server {
	driverName, err := dbdriver.SQLDriverName(cfg.DatabaseDriver)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open(driverName, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}
	if dbdriver.IsSQLite(cfg.DatabaseDriver) {
		conn.SetMaxOpenConns(1)
		if _, err := conn.Exec("PRAGMA foreign_keys = ON"); err != nil {
			log.Fatal("Cannot enable sqlite foreign keys:", err)
		}
	}

	dbQueries := database.New(conn)

	debug := cfg.Env != "production"
	router := NewRouter(dbQueries, conn, debug)
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
