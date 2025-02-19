package repository

import (
	"context"

	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/models"
)

type WordRepository interface {
	GetWord(ctx context.Context, id int) (*models.Word, error)
	ListWords(ctx context.Context, offset, limit int) ([]models.Word, error)
	CreateWord(ctx context.Context, word *models.Word) error
	UpdateWord(ctx context.Context, word *models.Word) error
	DeleteWord(ctx context.Context, id int) error
}

type GroupRepository interface {
	GetGroup(ctx context.Context, id int) (*models.Group, error)
	ListGroups(ctx context.Context, offset, limit int) ([]models.Group, error)
	CreateGroup(ctx context.Context, group *models.Group) error
	UpdateGroup(ctx context.Context, group *models.Group) error
	DeleteGroup(ctx context.Context, id int) error
}

type StudySessionRepository interface {
	GetStudySession(ctx context.Context, id int) (*models.StudySession, error)
	ListStudySessions(ctx context.Context, offset, limit int) ([]models.StudySession, error)
	CreateStudySession(ctx context.Context, session *models.StudySession) error
	AddWordReview(ctx context.Context, review *models.WordReviewItem) error
}
