package main

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/server"
	_ "modernc.org/sqlite"
)

func main() {
	cfg := config.Load()

	httpServer := server.NewServer(cfg)

	if err := httpServer.Start(); err != nil {
		log.Fatal(err)
	}
}
