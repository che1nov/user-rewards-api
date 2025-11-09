package usecases

import (
	"context"
	"fmt"

	"user-rewards-api/internal/dto"
)

type GetLeaderboardUseCase struct {
	postgres PostgreSQLAdapter
}

func NewGetLeaderboardUseCase(postgres PostgreSQLAdapter) *GetLeaderboardUseCase {
	return &GetLeaderboardUseCase{
		postgres: postgres,
	}
}

// Execute выполняет получение таблицы лидеров
func (uc *GetLeaderboardUseCase) Execute(ctx context.Context, limit int) (dto.GetLeaderboardOutput, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 100 {
		limit = 100
	}

	entries, err := uc.postgres.GetLeaderboard(ctx, limit)
	if err != nil {
		return dto.GetLeaderboardOutput{}, fmt.Errorf("ошибка при получении таблицы лидеров: %w", err)
	}

	leaderboardEntries := make([]dto.LeaderboardEntry, len(entries))
	for i := 0; i < len(entries); i++ {
		leaderboardEntries[i] = dto.LeaderboardEntry{
			Rank:     entries[i].Rank,
			UserID:   entries[i].UserID,
			Username: entries[i].Username,
			Balance:  entries[i].Balance,
		}
	}

	return dto.GetLeaderboardOutput{
		Users: leaderboardEntries,
		Total: len(leaderboardEntries),
	}, nil
}
