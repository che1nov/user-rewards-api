package usecases

import (
	"context"

	"user-rewards-api/internal/domain"
)

// LeaderboardEntry запись в таблице лидеров
type LeaderboardEntry struct {
	Rank     int
	UserID   string
	Username string
	Balance  int
}

// PostgreSQLAdapter интерфейс для работы с PostgreSQL
type PostgreSQLAdapter interface {
	// Методы для работы с пользователями
	CreateUser(ctx context.Context, user domain.User) error
	GetUserByID(ctx context.Context, userID domain.UserID) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUserBalance(ctx context.Context, userID domain.UserID, balance domain.Balance) error
	GetLeaderboard(ctx context.Context, limit int) ([]LeaderboardEntry, error)

	// Методы для работы с заданиями
	CreateTask(ctx context.Context, task domain.UserTask) error
	GetTasksByUserID(ctx context.Context, userID domain.UserID) ([]domain.UserTask, error)
	GetTaskByUserAndType(ctx context.Context, userID domain.UserID, taskType domain.TaskType) (*domain.UserTask, error)

	// Методы для работы с рефералами
	CreateReferral(ctx context.Context, referral domain.Referral) error
	GetReferralByReferredUserID(ctx context.Context, referredUserID domain.UserID) (*domain.Referral, error)
	CountReferralsByReferrerID(ctx context.Context, referrerID domain.UserID) (int, error)

	// Методы для работы с транзакциями
	WithTransaction(ctx context.Context, fn func(context.Context) error) error
}
