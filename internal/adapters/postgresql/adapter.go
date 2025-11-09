package postgresql

import (
	"context"

	"user-rewards-api/internal/domain"
	"user-rewards-api/internal/usecases"

	"github.com/jmoiron/sqlx"
)

// PostgreSQLAdapter объединяет все адаптеры PostgreSQL
type PostgreSQLAdapter struct {
	user        *PostgreSQLUserAdapter
	task        *PostgreSQLTaskAdapter
	referral    *PostgreSQLReferralAdapter
	transaction *PostgreSQLTransactionAdapter
}

// NewPostgreSQLAdapter создает новый объединенный адаптер PostgreSQL
func NewPostgreSQLAdapter(db *sqlx.DB) *PostgreSQLAdapter {
	return &PostgreSQLAdapter{
		user:        NewPostgreSQLUserAdapter(db),
		task:        NewPostgreSQLTaskAdapter(db),
		referral:    NewPostgreSQLReferralAdapter(db),
		transaction: NewPostgreSQLTransactionAdapter(db),
	}
}

// Методы для работы с пользователями
func (a *PostgreSQLAdapter) CreateUser(ctx context.Context, user domain.User) error {
	return a.user.CreateUser(ctx, user)
}

func (a *PostgreSQLAdapter) GetUserByID(ctx context.Context, userID domain.UserID) (*domain.User, error) {
	return a.user.GetUserByID(ctx, userID)
}

func (a *PostgreSQLAdapter) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	return a.user.GetUserByUsername(ctx, username)
}

func (a *PostgreSQLAdapter) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return a.user.GetUserByEmail(ctx, email)
}

func (a *PostgreSQLAdapter) UpdateUserBalance(ctx context.Context, userID domain.UserID, balance domain.Balance) error {
	return a.user.UpdateUserBalance(ctx, userID, balance)
}

func (a *PostgreSQLAdapter) GetLeaderboard(ctx context.Context, limit int) ([]usecases.LeaderboardEntry, error) {
	entries, err := a.user.GetLeaderboard(ctx, limit)
	if err != nil {
		return nil, err
	}

	result := make([]usecases.LeaderboardEntry, len(entries))
	for i := range entries {
		result[i] = usecases.LeaderboardEntry{
			Rank:     entries[i].Rank,
			UserID:   entries[i].UserID,
			Username: entries[i].Username,
			Balance:  entries[i].Balance,
		}
	}

	return result, nil
}

// Методы для работы с заданиями
func (a *PostgreSQLAdapter) CreateTask(ctx context.Context, task domain.UserTask) error {
	return a.task.CreateTask(ctx, task)
}

func (a *PostgreSQLAdapter) GetTasksByUserID(ctx context.Context, userID domain.UserID) ([]domain.UserTask, error) {
	return a.task.GetTasksByUserID(ctx, userID)
}

func (a *PostgreSQLAdapter) GetTaskByUserAndType(ctx context.Context, userID domain.UserID, taskType domain.TaskType) (*domain.UserTask, error) {
	return a.task.GetTaskByUserAndType(ctx, userID, taskType)
}

// Методы для работы с рефералами
func (a *PostgreSQLAdapter) CreateReferral(ctx context.Context, referral domain.Referral) error {
	return a.referral.CreateReferral(ctx, referral)
}

func (a *PostgreSQLAdapter) GetReferralByReferredUserID(ctx context.Context, referredUserID domain.UserID) (*domain.Referral, error) {
	return a.referral.GetReferralByReferredUserID(ctx, referredUserID)
}

func (a *PostgreSQLAdapter) CountReferralsByReferrerID(ctx context.Context, referrerID domain.UserID) (int, error) {
	return a.referral.CountReferralsByReferrerID(ctx, referrerID)
}

// Методы для работы с транзакциями
func (a *PostgreSQLAdapter) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	return a.transaction.WithTransaction(ctx, fn)
}
