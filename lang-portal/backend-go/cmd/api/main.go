package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/api/handlers"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/api/routes"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/repository/sqlite"
)

func main() {
	// Initialize SQLite database
	db, err := sql.Open("sqlite3", "words.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Apply migrations
	migrationSQL, err := os.ReadFile("database/migrations/001_initial_schema.sql")
	if err != nil {
		log.Fatal("Failed to read migration file:", err)
	}

	if _, err := db.Exec(string(migrationSQL)); err != nil {
		log.Fatal("Failed to apply migrations:", err)
	}

	// Initialize repositories
	wordRepo := sqlite.NewWordRepository(db)

	// Initialize handlers
	wordHandler := handlers.NewWordHandler(wordRepo)

	// Initialize Gin router
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r, wordHandler)

	// Start server
	log.Println("Server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
