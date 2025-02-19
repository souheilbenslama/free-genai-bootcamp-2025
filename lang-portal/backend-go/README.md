# Language Learning Portal Backend

A Go-based backend server for the language learning portal application.

## Features

- RESTful API using Gin framework
- SQLite database for data storage
- Supports vocabulary management with German-English translations
- Study session tracking
- Word grouping functionality

## Project Structure

```
backend-go/
├── cmd/api/              # Application entry point
├── internal/            # Private application code
├── database/           # Database migrations and seed data
├── config/            # Configuration files
├── pkg/               # Public packages
└── tests/            # Integration tests
```

## Prerequisites

- Go 1.16 or higher
- SQLite3

## Getting Started

1. Clone the repository
2. Navigate to the project directory
3. Install Mage (if not already installed):
   ```bash
   go install github.com/magefile/mage@latest
   ```
4. Install project dependencies:
   ```bash
   mage install
   ```
5. Initialize the database:
   ```bash
   mage db:init
   ```
6. Start the development server:
   ```bash
   mage dev:run
   ```

## Available Mage Commands

- `mage install` - Install project dependencies
- `mage build` - Build the application
- `mage test` - Run tests
- `mage clean` - Clean build artifacts

### Database Commands
- `mage db:init` - Initialize database with schema and seed data
- `mage db:migrate` - Apply database migrations
- `mage db:seed` - Load seed data
- `mage db:clean` - Remove database file

### Development Commands
- `mage dev:run` - Start development server

The server will start on `http://localhost:8080`

## API Endpoints

### Words

- `GET /api/words` - List all words
- `GET /api/words/:id` - Get a specific word
- `POST /api/words` - Create a new word

More endpoints coming soon.

## Development

The project uses:
- Gin for HTTP routing
- SQLite for data storage
- Standard Go project layout
