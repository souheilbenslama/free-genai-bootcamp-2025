package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/models"
)

type StudyRepository struct {
	db *sql.DB
}

func NewStudyRepository(db *sql.DB) *StudyRepository {
	return &StudyRepository{db: db}
}

func (r *StudyRepository) CreateStudySession(groupID int) (*models.StudySession, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", err)
	}
	defer tx.Rollback()

	// Create study session first
	sessionResult, err := tx.Exec(
		"INSERT INTO study_sessions (group_id, created_at) VALUES (?, ?)",
		groupID,
		time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating study session: %w", err)
	}

	sessionID, err := sessionResult.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting session ID: %w", err)
	}

	// Create study activity
	activityResult, err := tx.Exec(
		"INSERT INTO study_activities (study_session_id, group_id, created_at) VALUES (?, ?, ?)",
		sessionID,
		groupID,
		time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating study activity: %w", err)
	}

	activityID, err := activityResult.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting activity ID: %w", err)
	}

	// Update study session with activity ID
	_, err = tx.Exec(
		"UPDATE study_sessions SET study_activity_id = ? WHERE id = ?",
		activityID,
		sessionID,
	)
	if err != nil {
		return nil, fmt.Errorf("error updating study session: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &models.StudySession{
		ID:              int(sessionID),
		GroupID:         int(groupID),
		StudyActivityID: int(activityID),
		CreatedAt:       time.Now(),
	}, nil
}

func (r *StudyRepository) GetLastStudySession() (*models.StudySession, error) {
	query := `
		SELECT id, group_id, study_activity_id, created_at
		FROM study_sessions
		ORDER BY created_at DESC
		LIMIT 1
	`

	var session models.StudySession
	err := r.db.QueryRow(query).Scan(
		&session.ID,
		&session.GroupID,
		&session.StudyActivityID,
		&session.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error querying last study session: %w", err)
	}

	return &session, nil
}

func (r *StudyRepository) RecordWordReview(sessionID, wordID int, correct bool) error {
	_, err := r.db.Exec(
		"INSERT INTO word_review_items (word_id, study_session_id, correct, created_at) VALUES (?, ?, ?, ?)",
		wordID,
		sessionID,
		correct,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("error recording word review: %w", err)
	}

	return nil
}

func (r *StudyRepository) GetStudyProgress() (*models.StudyProgress, error) {
	// Get total available words
	var totalWords int
	err := r.db.QueryRow("SELECT COUNT(*) FROM words").Scan(&totalWords)
	if err != nil {
		return nil, fmt.Errorf("error counting words: %w", err)
	}

	// Get total words studied (unique words in word_review_items)
	var totalStudied int
	err = r.db.QueryRow(`
		SELECT COUNT(DISTINCT word_id)
		FROM word_review_items
	`).Scan(&totalStudied)
	if err != nil {
		return nil, fmt.Errorf("error counting studied words: %w", err)
	}

	// Calculate mastery percentage (words with > 80% correct answers)
	var masteredWords int
	err = r.db.QueryRow(`
		WITH word_stats AS (
			SELECT word_id,
				   COUNT(*) as total_reviews,
				   SUM(CASE WHEN correct THEN 1 ELSE 0 END) as correct_reviews
			FROM word_review_items
			GROUP BY word_id
		)
		SELECT COUNT(*)
		FROM word_stats
		WHERE CAST(correct_reviews AS FLOAT) / total_reviews >= 0.8
	`).Scan(&masteredWords)
	if err != nil {
		return nil, fmt.Errorf("error calculating mastery: %w", err)
	}

	masteryPercentage := 0.0
	if totalStudied > 0 {
		masteryPercentage = float64(masteredWords) / float64(totalStudied) * 100
	}

	return &models.StudyProgress{
		TotalWordsStudied:   totalStudied,
		TotalAvailableWords: totalWords,
		MasteryPercentage:   masteryPercentage,
	}, nil
}

func (r *StudyRepository) GetQuickStats() (*models.DashboardStats, error) {
	// Get success rate
	var correctCount, totalCount int
	err := r.db.QueryRow(`
		SELECT 
			COALESCE(SUM(CASE WHEN correct THEN 1 ELSE 0 END), 0) as correct_count,
			COALESCE(COUNT(*), 0) as total_count
		FROM word_review_items
	`).Scan(&correctCount, &totalCount)
	if err != nil {
		return nil, fmt.Errorf("error calculating success rate: %w", err)
	}

	successRate := 0.0
	if totalCount > 0 {
		successRate = float64(correctCount) / float64(totalCount)
	}

	// Get total study sessions
	var totalSessions int
	err = r.db.QueryRow("SELECT COALESCE(COUNT(*), 0) FROM study_sessions").Scan(&totalSessions)
	if err != nil {
		return nil, fmt.Errorf("error counting study sessions: %w", err)
	}

	// Get total active groups (groups with at least one study session)
	var activeGroups int
	err = r.db.QueryRow(`
		SELECT COALESCE(COUNT(DISTINCT group_id), 0)
		FROM study_sessions
	`).Scan(&activeGroups)
	if err != nil {
		return nil, fmt.Errorf("error counting active groups: %w", err)
	}

	// Calculate study streak (consecutive days with study sessions)
	var streak int
	err = r.db.QueryRow(`
		WITH RECURSIVE dates AS (
			SELECT date(created_at) as study_date
			FROM (
				SELECT created_at
				FROM study_sessions
				ORDER BY created_at DESC
				LIMIT 1
			)
		
			UNION ALL
		
			SELECT date(d.study_date, '-1 day')
			FROM dates d
			WHERE EXISTS (
				SELECT 1
				FROM study_sessions s
				WHERE date(s.created_at) = date(d.study_date, '-1 day')
			)
		)
		SELECT COALESCE(COUNT(*), 0)
		FROM dates
	`).Scan(&streak)
	if err != nil {
		return nil, fmt.Errorf("error calculating study streak: %w", err)
	}

	return &models.DashboardStats{
		SuccessRate:         successRate,
		TotalStudySessions: totalSessions,
		TotalActiveGroups:  activeGroups,
		StudyStreakDays:    streak,
	}, nil
}
