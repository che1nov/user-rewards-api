package usecases

import (
	"context"
	"fmt"

	"user-rewards-api/internal/domain"
	"user-rewards-api/internal/dto"
)

type ProcessReferralUseCase struct {
	postgres PostgreSQLAdapter
}

func NewProcessReferralUseCase(postgres PostgreSQLAdapter) *ProcessReferralUseCase {
	return &ProcessReferralUseCase{
		postgres: postgres,
	}
}

// Execute выполняет обработку реферального кода
func (uc *ProcessReferralUseCase) Execute(ctx context.Context, referredUserIDStr string, input dto.ProcessReferralInput) (dto.ProcessReferralOutput, error) {
	referredUserID, err := domain.UserIDFromString(referredUserIDStr)
	if err != nil {
		return dto.ProcessReferralOutput{}, err
	}

	referrerID, err := domain.UserIDFromString(input.ReferrerID)
	if err != nil {
		return dto.ProcessReferralOutput{}, domain.ErrReferrerNotFound
	}

	referredUser, err := uc.postgres.GetUserByID(ctx, referredUserID)
	if err != nil {
		return dto.ProcessReferralOutput{}, err
	}
	if referredUser == nil {
		return dto.ProcessReferralOutput{}, domain.ErrUserNotFound
	}

	referrer, err := uc.postgres.GetUserByID(ctx, referrerID)
	if err != nil {
		return dto.ProcessReferralOutput{}, err
	}
	if referrer == nil {
		return dto.ProcessReferralOutput{}, domain.ErrReferrerNotFound
	}

	existingReferral, err := uc.postgres.GetReferralByReferredUserID(ctx, referredUserID)
	if err != nil {
		return dto.ProcessReferralOutput{}, fmt.Errorf("ошибка при проверке реферальной связи: %w", err)
	}
	if existingReferral != nil {
		return dto.ProcessReferralOutput{}, domain.ErrReferralExists
	}

	var referral domain.Referral
	var newBalance domain.Balance

	err = uc.postgres.WithTransaction(ctx, func(ctx context.Context) error {
		referral, err = domain.NewReferral(referrerID, referredUserID)
		if err != nil {
			return err
		}

		if err := uc.postgres.CreateReferral(ctx, referral); err != nil {
			return fmt.Errorf("ошибка при создании реферальной связи: %w", err)
		}

		newBalance = referredUser.Balance.Add(referral.BonusPoints)
		if err := uc.postgres.UpdateUserBalance(ctx, referredUserID, newBalance); err != nil {
			return fmt.Errorf("ошибка при обновлении баланса: %w", err)
		}

		return nil
	})

	if err != nil {
		return dto.ProcessReferralOutput{}, err
	}

	return dto.ProcessReferralOutput{
		ReferralID:     referral.ID.String(),
		ReferrerID:     referrerID.String(),
		ReferredUserID: referredUserID.String(),
		BonusPoints:    referral.BonusPoints,
		NewBalance:     newBalance.Value(),
	}, nil
}
