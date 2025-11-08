package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UserID struct {
	value uuid.UUID
}

// NewUserID создает новый UserID
func NewUserID() (UserID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return UserID{}, fmt.Errorf("ошибка генерации ID: %w", err)
	}
	return UserID{value: id}, nil
}

// UserIDFromString создает UserID из строки
func UserIDFromString(s string) (UserID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return UserID{}, ErrInvalidUsername
	}
	return UserID{value: id}, nil
}

// String возвращает строковое представление UserID
func (id UserID) String() string {
	return id.value.String()
}

// Value возвращает UUID
func (id UserID) Value() uuid.UUID {
	return id.value
}

type Username struct {
	value string
}

// NewUsername создает новый Username с валидацией
func NewUsername(value string) (Username, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return Username{}, ErrInvalidUsername
	}
	if len(value) < 2 {
		return Username{}, fmt.Errorf("%w: минимальная длина 2 символа", ErrInvalidUsername)
	}
	if len(value) > 50 {
		return Username{}, fmt.Errorf("%w: максимальная длина 50 символов", ErrInvalidUsername)
	}
	return Username{value: value}, nil
}

// String возвращает строковое представление Username
func (u Username) String() string {
	return u.value
}

type Email struct {
	value string
}

// NewEmail создает новый Email с валидацией
func NewEmail(value string) (Email, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return Email{}, ErrInvalidEmail
	}
	if !strings.Contains(value, "@") {
		return Email{}, fmt.Errorf("%w: должен содержать @", ErrInvalidEmail)
	}
	if len(value) > 255 {
		return Email{}, fmt.Errorf("%w: максимальная длина 255 символов", ErrInvalidEmail)
	}
	return Email{value: value}, nil
}

// String возвращает строковое представление Email
func (e Email) String() string {
	return e.value
}

type Balance struct {
	value int
}

// NewBalance создает новый Balance
func NewBalance(value int) Balance {
	if value < 0 {
		value = 0
	}
	return Balance{value: value}
}

// Add добавляет поинты к балансу
func (b Balance) Add(points int) Balance {
	newValue := b.value + points
	if newValue < 0 {
		newValue = 0
	}
	return Balance{value: newValue}
}

// Value возвращает значение баланса
func (b Balance) Value() int {
	return b.value
}

type User struct {
	ID        UserID
	Username  Username
	Email     Email
	Balance   Balance
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser создает нового пользователя с валидацией
func NewUser(username, email string) (User, error) {
	userID, err := NewUserID()
	if err != nil {
		return User{}, err
	}

	usernameValue, err := NewUsername(username)
	if err != nil {
		return User{}, err
	}

	emailValue, err := NewEmail(email)
	if err != nil {
		return User{}, err
	}

	now := time.Now()
	return User{
		ID:        userID,
		Username:  usernameValue,
		Email:     emailValue,
		Balance:   NewBalance(0),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// AddPoints добавляет поинты к балансу пользователя
func (u *User) AddPoints(points int) {
	u.Balance = u.Balance.Add(points)
	u.UpdatedAt = time.Now()
}
