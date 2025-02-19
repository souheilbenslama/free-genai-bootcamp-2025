package models

import (
	"encoding/json"
	"time"
)

type Word struct {
	ID      int             `json:"id"`
	German  string          `json:"german"`
	English string          `json:"english"`
	Parts   json.RawMessage `json:"parts"`
}

type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
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
