package sqlite

import (
	"context"
	"database/sql"

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
	_, err := r.db.ExecContext(ctx,
		"UPDATE words SET german = ?, english = ?, parts = ? WHERE id = ?",
		word.German, word.English, word.Parts, word.ID)
	return err
}

func (r *WordRepository) DeleteWord(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM words WHERE id = ?", id)
	return err
}
