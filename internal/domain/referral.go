package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ReferralID struct {
	value uuid.UUID
}

// NewReferralID создает новый ReferralID
func NewReferralID() (ReferralID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return ReferralID{}, fmt.Errorf("ошибка генерации ID: %w", err)
	}
	return ReferralID{value: id}, nil
}

// ReferralIDFromString создает ReferralID из строки
func ReferralIDFromString(s string) (ReferralID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return ReferralID{}, fmt.Errorf("некорректный формат ReferralID: %w", err)
	}
	return ReferralID{value: id}, nil
}

// String возвращает строковое представление ReferralID
func (id ReferralID) String() string {
	return id.value.String()
}

// Value возвращает UUID
func (id ReferralID) Value() uuid.UUID {
	return id.value
}

// Referral представляет реферальную связь между пользователями
type Referral struct {
	ID             ReferralID
	ReferrerID     UserID
	ReferredUserID UserID
	BonusPoints    int
	CreatedAt      time.Time
}

// NewReferral создает новую реферальную связь
func NewReferral(referrerID, referredUserID UserID) (Referral, error) {
	if referrerID.String() == referredUserID.String() {
		return Referral{}, ErrSelfReferral
	}

	referralID, err := NewReferralID()
	if err != nil {
		return Referral{}, err
	}

	return Referral{
		ID:             referralID,
		ReferrerID:     referrerID,
		ReferredUserID: referredUserID,
		BonusPoints:    100,
		CreatedAt:      time.Now(),
	}, nil
}
