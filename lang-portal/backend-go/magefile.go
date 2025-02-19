//go:build mage
package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	_ "github.com/mattn/go-sqlite3"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/seeder"
)

type DB mg.Namespace

// Init initializes the database with schema and seed data
func (DB) Init() error {
	mg.Deps(DB.Migrate)
	return DB{}.Seed()
}

// Migrate applies database migrations
func (DB) Migrate() error {
	fmt.Println("Applying database migrations...")
	
	db, err := sql.Open("sqlite3", "words.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	migrationPath := filepath.Join("database", "migrations", "001_initial_schema.sql")
	migrationSQL, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	if _, err := db.Exec(string(migrationSQL)); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	fmt.Println("Migrations applied successfully")
	return nil
}

// Seed loads initial data into the database
func (DB) Seed() error {
	fmt.Println("Loading seed data...")
	
	db, err := sql.Open("sqlite3", "words.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	seedDir := filepath.Join("database", "seed")
	if err := seeder.LoadSeedData(db, seedDir); err != nil {
		return fmt.Errorf("failed to load seed data: %w", err)
	}

	fmt.Println("Seed data loaded successfully")
	return nil
}

// Clean removes the database file
func (DB) Clean() error {
	fmt.Println("Cleaning database...")
	return os.Remove("words.db")
}

type Dev mg.Namespace

// Run starts the development server
func (Dev) Run() error {
	mg.Deps(DB.Init)
	fmt.Println("Starting development server...")
	return sh.Run("go", "run", "./cmd/api")
}

// Build builds the application
func Build() error {
	fmt.Println("Building application...")
	return sh.Run("go", "build", "-o", "lang-portal", "./cmd/api")
}

// Test runs the test suite
func Test() error {
	fmt.Println("Running tests...")
	return sh.Run("go", "test", "./...")
}

// Clean removes build artifacts
func Clean() error {
	fmt.Println("Cleaning build artifacts...")
	mg.Deps(DB.Clean)
	return os.Remove("lang-portal")
}

// Install installs project dependencies
func Install() error {
	fmt.Println("Installing dependencies...")
	return sh.Run("go", "mod", "download")
}
