package usecases

import (
	"context"
	"fmt"

	"user-rewards-api/internal/domain"
	"user-rewards-api/internal/dto"
)

type CompleteTaskUseCase struct {
	postgres PostgreSQLAdapter
}

func NewCompleteTaskUseCase(postgres PostgreSQLAdapter) *CompleteTaskUseCase {
	return &CompleteTaskUseCase{
		postgres: postgres,
	}
}

// Execute выполняет задание пользователя
func (uc *CompleteTaskUseCase) Execute(ctx context.Context, userIDStr string, input dto.CompleteTaskInput) (dto.CompleteTaskOutput, error) {
	userID, err := domain.UserIDFromString(userIDStr)
	if err != nil {
		return dto.CompleteTaskOutput{}, err
	}

	taskType, err := domain.NewTaskType(input.TaskType)
	if err != nil {
		return dto.CompleteTaskOutput{}, err
	}

	user, err := uc.postgres.GetUserByID(ctx, userID)
	if err != nil {
		return dto.CompleteTaskOutput{}, err
	}
	if user == nil {
		return dto.CompleteTaskOutput{}, domain.ErrUserNotFound
	}

	existingTask, err := uc.postgres.GetTaskByUserAndType(ctx, userID, taskType)
	if err != nil {
		return dto.CompleteTaskOutput{}, fmt.Errorf("ошибка при проверке задания: %w", err)
	}
	if existingTask != nil {
		return dto.CompleteTaskOutput{}, domain.ErrTaskAlreadyExists
	}

	var task domain.UserTask
	var newBalance domain.Balance

	err = uc.postgres.WithTransaction(ctx, func(ctx context.Context) error {
		task, err = domain.NewUserTask(userID, taskType)
		if err != nil {
			return err
		}

		if err := uc.postgres.CreateTask(ctx, task); err != nil {
			return fmt.Errorf("ошибка при создании задания: %w", err)
		}

		newBalance = user.Balance.Add(taskType.GetPoints())
		if err := uc.postgres.UpdateUserBalance(ctx, userID, newBalance); err != nil {
			return fmt.Errorf("ошибка при обновлении баланса: %w", err)
		}

		if taskType == domain.TaskTypeInviteFriend {
			referral, err := uc.postgres.GetReferralByReferredUserID(ctx, userID)
			if err != nil {
				return fmt.Errorf("ошибка при проверке реферальной связи: %w", err)
			}

			if referral != nil {
				referrer, err := uc.postgres.GetUserByID(ctx, referral.ReferrerID)
				if err != nil {
					return fmt.Errorf("ошибка при получении реферера: %w", err)
				}
				if referrer != nil {
					referrerNewBalance := referrer.Balance.Add(taskType.GetPoints())
					if err := uc.postgres.UpdateUserBalance(ctx, referral.ReferrerID, referrerNewBalance); err != nil {
						return fmt.Errorf("ошибка при обновлении баланса реферера: %w", err)
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return dto.CompleteTaskOutput{}, err
	}

	return dto.CompleteTaskOutput{
		TaskID:     task.ID.String(),
		TaskType:   taskType.String(),
		Points:     taskType.GetPoints(),
		NewBalance: newBalance.Value(),
	}, nil
}
