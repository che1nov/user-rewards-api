package usecases

import (
	"context"
	"fmt"

	"user-rewards-api/internal/domain"
	"user-rewards-api/internal/dto"
)

type GetUserStatusUseCase struct {
	postgres PostgreSQLAdapter
}

func NewGetUserStatusUseCase(postgres PostgreSQLAdapter) *GetUserStatusUseCase {
	return &GetUserStatusUseCase{
		postgres: postgres,
	}
}

// Execute выполняет получение статуса пользователя
func (uc *GetUserStatusUseCase) Execute(ctx context.Context, userIDStr string) (dto.GetUserStatusOutput, error) {
	userID, err := domain.UserIDFromString(userIDStr)
	if err != nil {
		return dto.GetUserStatusOutput{}, err
	}

	user, err := uc.postgres.GetUserByID(ctx, userID)
	if err != nil {
		return dto.GetUserStatusOutput{}, err
	}
	if user == nil {
		return dto.GetUserStatusOutput{}, domain.ErrUserNotFound
	}

	tasks, err := uc.postgres.GetTasksByUserID(ctx, userID)
	if err != nil {
		return dto.GetUserStatusOutput{}, fmt.Errorf("ошибка при получении заданий: %w", err)
	}

	referralCount, err := uc.postgres.CountReferralsByReferrerID(ctx, userID)
	if err != nil {
		return dto.GetUserStatusOutput{}, fmt.Errorf("ошибка при получении рефералов: %w", err)
	}

	return dto.GetUserStatusOutput{
		UserID:         user.ID.String(),
		Balance:        user.Balance.Value(),
		CompletedTasks: len(tasks),
		ReferralCount:  referralCount,
	}, nil
}
