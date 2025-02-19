package models

import "time"

type Word struct {
	ID      int    `json:"id"`
	German  string `json:"german"`
	English string `json:"english"`
	Parts   string `json:"parts"`
}

type Group struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	WordCount   int       `json:"word_count,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type StudySession struct {
	ID              int       `json:"id"`
	GroupID         int       `json:"group_id"`
	CreatedAt       time.Time `json:"created_at"`
	StudyActivityID int       `json:"study_activity_id"`
}

type StudyActivity struct {
	ID              int       `json:"id"`
	StudySessionID  int       `json:"study_session_id"`
	GroupID         int       `json:"group_id"`
	CreatedAt       time.Time `json:"created_at"`
}

type WordReviewItem struct {
	WordID         int       `json:"word_id"`
	StudySessionID int       `json:"study_session_id"`
	Correct        bool      `json:"correct"`
	CreatedAt      time.Time `json:"created_at"`
}

type DashboardStats struct {
	SuccessRate         float64 `json:"success_rate"`
	TotalStudySessions int     `json:"total_study_sessions"`
	TotalActiveGroups  int     `json:"total_active_groups"`
	StudyStreakDays    int     `json:"study_streak_days"`
}

type StudyProgress struct {
	TotalWordsStudied    int     `json:"total_words_studied"`
	TotalAvailableWords  int     `json:"total_available_words"`
	MasteryPercentage    float64 `json:"mastery_percentage"`
}
