package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TaskType string

const (
	TaskTypeSurvey            TaskType = "survey"
	TaskTypeSubscribeTelegram TaskType = "subscribe_telegram"
	TaskTypeSubscribeTwitter  TaskType = "subscribe_twitter"
	TaskTypeInviteFriend      TaskType = "invite_friend"
)

// NewTaskType создает новый TaskType с валидацией
func NewTaskType(value string) (TaskType, error) {
	taskType := TaskType(value)
	if !taskType.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidTaskType, value)
	}
	return taskType, nil
}

// GetPoints возвращает количество поинтов за выполнение задания
func (t TaskType) GetPoints() int {
	switch t {
	case TaskTypeSurvey:
		return 10
	case TaskTypeSubscribeTelegram:
		return 50
	case TaskTypeSubscribeTwitter:
		return 50
	case TaskTypeInviteFriend:
		return 100
	default:
		return 0
	}
}

// IsValid проверяет валидность типа задания
func (t TaskType) IsValid() bool {
	return t == TaskTypeSurvey ||
		t == TaskTypeSubscribeTelegram ||
		t == TaskTypeSubscribeTwitter ||
		t == TaskTypeInviteFriend
}

// String возвращает строковое представление TaskType
func (t TaskType) String() string {
	return string(t)
}

// TaskID представляет идентификатор задания
type TaskID struct {
	value uuid.UUID
}

// NewTaskID создает новый TaskID
func NewTaskID() (TaskID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return TaskID{}, fmt.Errorf("ошибка генерации ID: %w", err)
	}
	return TaskID{value: id}, nil
}

// TaskIDFromString создает TaskID из строки
func TaskIDFromString(s string) (TaskID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return TaskID{}, fmt.Errorf("некорректный формат TaskID: %w", err)
	}
	return TaskID{value: id}, nil
}

// String возвращает строковое представление TaskID
func (id TaskID) String() string {
	return id.value.String()
}

// Value возвращает UUID
func (id TaskID) Value() uuid.UUID {
	return id.value
}

// UserTask представляет выполненное задание пользователя
type UserTask struct {
	ID          TaskID
	UserID      UserID
	TaskType    TaskType
	CompletedAt time.Time
	Points      int
}

// NewUserTask создает новое задание пользователя
func NewUserTask(userID UserID, taskType TaskType) (UserTask, error) {
	taskID, err := NewTaskID()
	if err != nil {
		return UserTask{}, err
	}

	return UserTask{
		ID:          taskID,
		UserID:      userID,
		TaskType:    taskType,
		CompletedAt: time.Now(),
		Points:      taskType.GetPoints(),
	}, nil
}
