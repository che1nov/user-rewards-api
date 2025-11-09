package postgresql

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"user-rewards-api/internal/domain"
)

// PostgreSQLReferralAdapter адаптер для работы с рефералами в PostgreSQL
type PostgreSQLReferralAdapter struct {
	db *sqlx.DB
}

// NewPostgreSQLReferralAdapter создает новый адаптер рефералов
func NewPostgreSQLReferralAdapter(db *sqlx.DB) *PostgreSQLReferralAdapter {
	return &PostgreSQLReferralAdapter{db: db}
}

// CreateReferral создает новую реферальную связь
func (a *PostgreSQLReferralAdapter) CreateReferral(ctx context.Context, referral domain.Referral) error {
	query := `
		INSERT INTO referrals (id, referrer_id, referred_user_id, bonus_points, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := a.db.ExecContext(ctx, query,
		referral.ID.Value(), referral.ReferrerID.Value(), referral.ReferredUserID.Value(),
		referral.BonusPoints, referral.CreatedAt)
	return err
}

// GetReferralByReferredUserID получает реферальную связь по ID приглашенного пользователя
func (a *PostgreSQLReferralAdapter) GetReferralByReferredUserID(ctx context.Context, referredUserID domain.UserID) (*domain.Referral, error) {
	var referral struct {
		ID             string    `db:"id"`
		ReferrerID     string    `db:"referrer_id"`
		ReferredUserID string    `db:"referred_user_id"`
		BonusPoints    int       `db:"bonus_points"`
		CreatedAt      time.Time `db:"created_at"`
	}

	query := `
		SELECT id, referrer_id, referred_user_id, bonus_points, created_at
		FROM referrals
		WHERE referred_user_id = $1
	`

	err := a.db.GetContext(ctx, &referral, query, referredUserID.Value())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	referralID, err := domain.ReferralIDFromString(referral.ID)
	if err != nil {
		return nil, err
	}

	referrerID, err := domain.UserIDFromString(referral.ReferrerID)
	if err != nil {
		return nil, err
	}

	referredID, err := domain.UserIDFromString(referral.ReferredUserID)
	if err != nil {
		return nil, err
	}

	return &domain.Referral{
		ID:             referralID,
		ReferrerID:     referrerID,
		ReferredUserID: referredID,
		BonusPoints:    referral.BonusPoints,
		CreatedAt:      referral.CreatedAt,
	}, nil
}

// CountReferralsByReferrerID подсчитывает количество рефералов по ID реферера
func (a *PostgreSQLReferralAdapter) CountReferralsByReferrerID(ctx context.Context, referrerID domain.UserID) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM referrals WHERE referrer_id = $1`

	err := a.db.GetContext(ctx, &count, query, referrerID.Value())
	if err != nil {
		return 0, err
	}
	return count, nil
}

