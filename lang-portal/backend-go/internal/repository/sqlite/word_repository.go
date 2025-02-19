package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/models"
)

type WordRepository struct {
	db *sql.DB
}



func NewWordRepository(db *sql.DB) *WordRepository {
	return &WordRepository{db: db}
}

func (r *WordRepository) GetWord(ctx context.Context, id int) (*models.Word, error) {
	var word models.Word
	err := r.db.QueryRowContext(ctx,
		"SELECT id, german, english, parts FROM words WHERE id = ?",
		id).Scan(&word.ID, &word.German, &word.English, &word.Parts)
	if err != nil {
		return nil, err
	}
	return &word, nil
}

func (r *WordRepository) ListWords(ctx context.Context, offset, limit int) ([]models.Word, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, german, english, parts FROM words LIMIT ? OFFSET ?",
		limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []models.Word
	for rows.Next() {
		var word models.Word
		if err := rows.Scan(&word.ID, &word.German, &word.English, &word.Parts); err != nil {
			return nil, err
		}
		words = append(words, word)
	}
	return words, nil
}

func (r *WordRepository) CreateWord(ctx context.Context, word *models.Word) error {
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO words (german, english, parts) VALUES (?, ?, ?)",
		word.German, word.English, word.Parts)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	word.ID = int(id)
	return nil
}

func (r *WordRepository) UpdateWord(ctx context.Context, word *models.Word) error {
	query := `
		UPDATE words
		SET german = ?, english = ?, parts = ?
		WHERE id = ?
		RETURNING id, german, english, parts
	`

	err := r.db.QueryRowContext(ctx,
		query,
		word.German,
		word.English,
		word.Parts,
		word.ID,
	).Scan(
		&word.ID,
		&word.German,
		&word.English,
		&word.Parts,
	)

	if err != nil {
		return fmt.Errorf("error updating word: %w", err)
	}

	return nil
}

func (r *WordRepository) DeleteWord(ctx context.Context, id int) error {
	// First delete any associations in words_groups
	_, err := r.db.ExecContext(ctx, "DELETE FROM words_groups WHERE word_id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting word associations: %w", err)
	}

	// Then delete the word
	result, err := r.db.ExecContext(ctx, "DELETE FROM words WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting word: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("word with id %d not found", id)
	}

	return nil
}
