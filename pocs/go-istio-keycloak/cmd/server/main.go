package main

import (
	"log"

	"book-service/internal/config"
	"book-service/internal/db"
	"book-service/internal/repository"
	"book-service/internal/router"
)

func main() {
	cfg := config.Load()

	conn, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer conn.Close()

	if err := db.Migrate(conn); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	repo := repository.NewBookRepository(conn)
	r := router.New(cfg, repo)

	log.Printf("book-service listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("run: %v", err)
	}
}
