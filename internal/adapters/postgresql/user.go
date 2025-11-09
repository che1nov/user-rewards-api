package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"user-rewards-api/internal/domain"

	"github.com/jmoiron/sqlx"
)

type PostgreSQLUserAdapter struct {
	db *sqlx.DB
}

func NewPostgreSQLUserAdapter(db *sqlx.DB) *PostgreSQLUserAdapter {
	return &PostgreSQLUserAdapter{db: db}
}

// CreateUser создает нового пользователя
func (a *PostgreSQLUserAdapter) CreateUser(ctx context.Context, user domain.User) error {
	query := `
		INSERT INTO users (id, username, email, balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := a.db.ExecContext(ctx, query,
		user.ID.Value(), user.Username.String(), user.Email.String(),
		user.Balance.Value(), user.CreatedAt, user.UpdatedAt)
	return err
}

// GetUserByID получает пользователя по ID
func (a *PostgreSQLUserAdapter) GetUserByID(ctx context.Context, userID domain.UserID) (*domain.User, error) {
	var user struct {
		ID        string    `db:"id"`
		Username  string    `db:"username"`
		Email     string    `db:"email"`
		Balance   int       `db:"balance"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	query := `SELECT id, username, email, balance, created_at, updated_at FROM users WHERE id = $1`
	err := a.db.GetContext(ctx, &user, query, userID.Value())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	domainUserID, err := domain.UserIDFromString(user.ID)
	if err != nil {
		return nil, err
	}

	username, err := domain.NewUsername(user.Username)
	if err != nil {
		return nil, err
	}

	email, err := domain.NewEmail(user.Email)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        domainUserID,
		Username:  username,
		Email:     email,
		Balance:   domain.NewBalance(user.Balance),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetUserByUsername получает пользователя по username
func (a *PostgreSQLUserAdapter) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user struct {
		ID        string    `db:"id"`
		Username  string    `db:"username"`
		Email     string    `db:"email"`
		Balance   int       `db:"balance"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	query := `SELECT id, username, email, balance, created_at, updated_at FROM users WHERE username = $1`
	err := a.db.GetContext(ctx, &user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	domainUserID, err := domain.UserIDFromString(user.ID)
	if err != nil {
		return nil, err
	}

	usernameValue, err := domain.NewUsername(user.Username)
	if err != nil {
		return nil, err
	}

	email, err := domain.NewEmail(user.Email)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        domainUserID,
		Username:  usernameValue,
		Email:     email,
		Balance:   domain.NewBalance(user.Balance),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetUserByEmail получает пользователя по email
func (a *PostgreSQLUserAdapter) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user struct {
		ID        string    `db:"id"`
		Username  string    `db:"username"`
		Email     string    `db:"email"`
		Balance   int       `db:"balance"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	query := `SELECT id, username, email, balance, created_at, updated_at FROM users WHERE email = $1`
	err := a.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	domainUserID, err := domain.UserIDFromString(user.ID)
	if err != nil {
		return nil, err
	}

	username, err := domain.NewUsername(user.Username)
	if err != nil {
		return nil, err
	}

	emailValue, err := domain.NewEmail(user.Email)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        domainUserID,
		Username:  username,
		Email:     emailValue,
		Balance:   domain.NewBalance(user.Balance),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// UpdateUserBalance обновляет баланс пользователя
func (a *PostgreSQLUserAdapter) UpdateUserBalance(ctx context.Context, userID domain.UserID, balance domain.Balance) error {
	query := `UPDATE users SET balance = $1, updated_at = $2 WHERE id = $3`
	_, err := a.db.ExecContext(ctx, query, balance.Value(), time.Now(), userID.Value())
	return err
}

// leaderboardRow представляет строку результата запроса leaderboard
type leaderboardRow struct {
	UserID   string `db:"user_id"`
	Username string `db:"username"`
	Balance  int    `db:"balance"`
}

// leaderboardEntry представляет запись в таблице лидеров (локальный тип для адаптера)
type leaderboardEntry struct {
	Rank     int
	UserID   string
	Username string
	Balance  int
}

// GetLeaderboard получает топ пользователей по балансу
func (a *PostgreSQLUserAdapter) GetLeaderboard(ctx context.Context, limit int) ([]leaderboardEntry, error) {
	query := `
		SELECT 
			id::text as user_id,
			username,
			balance
		FROM users
		ORDER BY balance DESC, created_at ASC
		LIMIT $1
	`

	var rows []leaderboardRow
	err := a.db.SelectContext(ctx, &rows, query, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения SQL запроса leaderboard: %w", err)
	}

	result := make([]leaderboardEntry, len(rows))
	for i := range rows {
		result[i] = leaderboardEntry{
			Rank:     i + 1,
			UserID:   rows[i].UserID,
			Username: rows[i].Username,
			Balance:  rows[i].Balance,
		}
	}

	return result, nil
}
