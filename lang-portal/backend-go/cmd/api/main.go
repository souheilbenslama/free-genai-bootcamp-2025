package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/api/handlers"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/api/routes"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/repository/sqlite"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/seeder"
)

// getProjectRoot returns the absolute path to the project root directory
func getProjectRoot() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}
	// Go up two directories from cmd/api/main.go to reach project root
	projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")
	return filepath.Abs(projectRoot)
}

func main() {
	// Get the project root directory
	projectRoot, err := getProjectRoot()
	if err != nil {
		log.Fatal("Failed to get project root:", err)
	}

	// Initialize SQLite database
	dbPath := filepath.Join(projectRoot, "words.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Apply migrations
	migrationPath := filepath.Join(projectRoot, "database", "migrations", "001_initial_schema.sql")
	log.Printf("Looking for migration file at: %s", migrationPath)
	
	migrationSQL, err := os.ReadFile(migrationPath)
	if err != nil {
		log.Fatal("Failed to read migration file:", err)
	}

	if _, err := db.Exec(string(migrationSQL)); err != nil {
		log.Fatal("Failed to apply migrations:", err)
	}

	// Load seed data from JSON files
	seedDir := filepath.Join(projectRoot, "database", "seed")
	log.Printf("Loading seed data from directory: %s", seedDir)
	
	if err := seeder.LoadSeedData(db, seedDir); err != nil {
		log.Fatal("Failed to load seed data:", err)
	}

	log.Println("Database initialized with seed data")

	// Initialize repositories
	wordRepo := sqlite.NewWordRepository(db)
	groupRepo := sqlite.NewGroupRepository(db)
	studyRepo := sqlite.NewStudyRepository(db)

	// Initialize handlers
	wordHandler := handlers.NewWordHandler(wordRepo)
	groupHandler := handlers.NewGroupHandler(groupRepo)
	studyHandler := handlers.NewStudyHandler(studyRepo)

	// Initialize Gin router
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r, wordHandler, groupHandler, studyHandler)

	// Start server
	log.Printf("Server starting on :8080... (Project root: %s)", projectRoot)
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
