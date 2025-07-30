package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/othavioBF/pandoragym-go-api/internal/core"
	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	gob.Register(uuid.UUID{})

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	pool, err := pgstore.InitDB(databaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer pool.Close()

	queries := pgstore.NewQueries(pool)

	// Initialize API with dependency injection
	api := core.InjectDependencies(queries, pool)

	// Bind routes
	api.BindRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3333"
	}

	fmt.Printf("🚀 PandoraGym API Server starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, api.Router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
