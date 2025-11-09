package usecases

import (
	"context"
	"fmt"
	"time"

	"user-rewards-api/internal/domain"
	"user-rewards-api/internal/dto"

	"github.com/golang-jwt/jwt/v5"
)

type CreateUserUseCase struct {
	postgres  PostgreSQLAdapter
	jwtSecret string
}

func NewCreateUserUseCase(postgres PostgreSQLAdapter, jwtSecret string) *CreateUserUseCase {
	return &CreateUserUseCase{
		postgres:  postgres,
		jwtSecret: jwtSecret,
	}
}

// Execute выполняет создание пользователя
func (uc *CreateUserUseCase) Execute(ctx context.Context, input dto.CreateUserInput) (dto.CreateUserOutput, error) {
	if existingUser, err := uc.postgres.GetUserByUsername(ctx, input.Username); err != nil {
		return dto.CreateUserOutput{}, fmt.Errorf("ошибка при проверке username: %w", err)
	} else if existingUser != nil {
		return dto.CreateUserOutput{}, domain.ErrUserExists
	}

	if existingUser, err := uc.postgres.GetUserByEmail(ctx, input.Email); err != nil {
		return dto.CreateUserOutput{}, fmt.Errorf("ошибка при проверке email: %w", err)
	} else if existingUser != nil {
		return dto.CreateUserOutput{}, domain.ErrUserExists
	}

	user, err := domain.NewUser(input.Username, input.Email)
	if err != nil {
		return dto.CreateUserOutput{}, err
	}

	if err := uc.postgres.CreateUser(ctx, user); err != nil {
		return dto.CreateUserOutput{}, fmt.Errorf("ошибка при создании пользователя: %w", err)
	}

	token, err := uc.generateJWT(user.ID.String())
	if err != nil {
		return dto.CreateUserOutput{}, fmt.Errorf("ошибка при генерации токена: %w", err)
	}

	return dto.CreateUserOutput{
		UserID:      user.ID.String(),
		Username:    user.Username.String(),
		Email:       user.Email.String(),
		AccessToken: token,
	}, nil
}

// generateJWT генерирует JWT токен для пользователя
func (uc *CreateUserUseCase) generateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}
