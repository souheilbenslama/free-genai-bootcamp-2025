package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/models"
)

type GroupRepository struct {
	db *sql.DB
}

func (r *GroupRepository) CreateGroup(ctx context.Context, group *models.Group) error {
	query := `
		INSERT INTO groups (name, description)
		VALUES (?, ?)
		RETURNING id, name, description
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		group.Name,
		group.Description,
	).Scan(
		&group.ID,
		&group.Name,
		&group.Description,
	)

	if err != nil {
		return fmt.Errorf("error creating group: %w", err)
	}

	return nil
}

func (r *GroupRepository) UpdateGroup(ctx context.Context, group *models.Group) error {
	query := `
		UPDATE groups
		SET name = ?, description = ?
		WHERE id = ?
		RETURNING id, name, description
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		group.Name,
		group.Description,
		group.ID,
	).Scan(
		&group.ID,
		&group.Name,
		&group.Description,
	)

	if err != nil {
		return fmt.Errorf("error updating group: %w", err)
	}

	return nil
}

func (r *GroupRepository) DeleteGroup(ctx context.Context, id int) error {
	// First delete any word associations
	_, err := r.db.ExecContext(ctx, "DELETE FROM words_groups WHERE group_id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting group associations: %w", err)
	}

	// Then delete the group
	result, err := r.db.ExecContext(ctx, "DELETE FROM groups WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting group: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("group with id %d not found", id)
	}

	return nil
}

func (r *GroupRepository) AddWordToGroup(ctx context.Context, groupID, wordID int) error {
	// Check if the association already exists
	var exists bool
	err := r.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM words_groups WHERE group_id = ? AND word_id = ?)",
		groupID,
		wordID,
	).Scan(&exists)

	if err != nil {
		return fmt.Errorf("error checking word-group association: %w", err)
	}

	if exists {
		return fmt.Errorf("word %d is already in group %d", wordID, groupID)
	}

	// Add the association
	_, err = r.db.ExecContext(
		ctx,
		"INSERT INTO words_groups (group_id, word_id) VALUES (?, ?)",
		groupID,
		wordID,
	)

	if err != nil {
		return fmt.Errorf("error adding word to group: %w", err)
	}

	return nil
}

func (r *GroupRepository) RemoveWordFromGroup(ctx context.Context, groupID, wordID int) error {
	result, err := r.db.ExecContext(
		ctx,
		"DELETE FROM words_groups WHERE group_id = ? AND word_id = ?",
		groupID,
		wordID,
	)

	if err != nil {
		return fmt.Errorf("error removing word from group: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("word %d is not in group %d", wordID, groupID)
	}

	return nil
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) GetAll(page, pageSize int) ([]models.Group, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM groups").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting groups: %w", err)
	}

	// Get groups with word count
	query := `
		SELECT g.id, g.name, COUNT(wg.word_id) as word_count
		FROM groups g
		LEFT JOIN words_groups wg ON g.id = wg.group_id
		GROUP BY g.id
		ORDER BY g.id
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying groups: %w", err)
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var g models.Group
		if err := rows.Scan(&g.ID, &g.Name, &g.WordCount); err != nil {
			return nil, 0, fmt.Errorf("error scanning group: %w", err)
		}
		groups = append(groups, g)
	}

	return groups, total, nil
}

func (r *GroupRepository) GetByID(id int) (*models.Group, error) {
	query := `
		SELECT g.id, g.name, COUNT(wg.word_id) as word_count
		FROM groups g
		LEFT JOIN words_groups wg ON g.id = wg.group_id
		WHERE g.id = ?
		GROUP BY g.id
	`

	var group models.Group
	err := r.db.QueryRow(query, id).Scan(&group.ID, &group.Name, &group.WordCount)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error querying group: %w", err)
	}

	return &group, nil
}

func (r *GroupRepository) GetGroupWords(groupID int) ([]models.Word, error) {
	query := `
		SELECT w.id, w.german, w.english, w.parts
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?
		ORDER BY w.id
	`

	rows, err := r.db.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("error querying group words: %w", err)
	}
	defer rows.Close()

	var words []models.Word
	for rows.Next() {
		var w models.Word
		if err := rows.Scan(&w.ID, &w.German, &w.English, &w.Parts); err != nil {
			return nil, fmt.Errorf("error scanning word: %w", err)
		}
		words = append(words, w)
	}

	return words, nil
}
