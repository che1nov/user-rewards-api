package postgresql

import (
	"context"
	"database/sql"
	"time"

	"user-rewards-api/internal/domain"

	"github.com/jmoiron/sqlx"
)

type PostgreSQLTaskAdapter struct {
	db *sqlx.DB
}

func NewPostgreSQLTaskAdapter(db *sqlx.DB) *PostgreSQLTaskAdapter {
	return &PostgreSQLTaskAdapter{db: db}
}

// CreateTask создает новое задание
func (a *PostgreSQLTaskAdapter) CreateTask(ctx context.Context, task domain.UserTask) error {
	query := `
		INSERT INTO user_tasks (id, user_id, task_type, completed_at, points)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := a.db.ExecContext(ctx, query,
		task.ID.Value(), task.UserID.Value(), task.TaskType.String(),
		task.CompletedAt, task.Points)
	return err
}

// GetTasksByUserID получает все задания пользователя
func (a *PostgreSQLTaskAdapter) GetTasksByUserID(ctx context.Context, userID domain.UserID) ([]domain.UserTask, error) {
	var tasks []struct {
		ID          string    `db:"id"`
		UserID      string    `db:"user_id"`
		TaskType    string    `db:"task_type"`
		CompletedAt time.Time `db:"completed_at"`
		Points      int       `db:"points"`
	}

	query := `
		SELECT id, user_id, task_type, completed_at, points 
		FROM user_tasks 
		WHERE user_id = $1
		ORDER BY completed_at DESC
	`

	err := a.db.SelectContext(ctx, &tasks, query, userID.Value())
	if err != nil {
		return nil, err
	}

	result := make([]domain.UserTask, 0, len(tasks))
	for _, t := range tasks {
		taskID, err := domain.TaskIDFromString(t.ID)
		if err != nil {
			return nil, err
		}

		domainUserID, err := domain.UserIDFromString(t.UserID)
		if err != nil {
			return nil, err
		}

		taskType, err := domain.NewTaskType(t.TaskType)
		if err != nil {
			return nil, err
		}

		result = append(result, domain.UserTask{
			ID:          taskID,
			UserID:      domainUserID,
			TaskType:    taskType,
			CompletedAt: t.CompletedAt,
			Points:      t.Points,
		})
	}

	return result, nil
}

// GetTaskByUserAndType получает задание пользователя по типу
func (a *PostgreSQLTaskAdapter) GetTaskByUserAndType(ctx context.Context, userID domain.UserID, taskType domain.TaskType) (*domain.UserTask, error) {
	var task struct {
		ID          string    `db:"id"`
		UserID      string    `db:"user_id"`
		TaskType    string    `db:"task_type"`
		CompletedAt time.Time `db:"completed_at"`
		Points      int       `db:"points"`
	}

	query := `
		SELECT id, user_id, task_type, completed_at, points 
		FROM user_tasks 
		WHERE user_id = $1 AND task_type = $2
	`

	err := a.db.GetContext(ctx, &task, query, userID.Value(), taskType.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	taskID, err := domain.TaskIDFromString(task.ID)
	if err != nil {
		return nil, err
	}

	domainUserID, err := domain.UserIDFromString(task.UserID)
	if err != nil {
		return nil, err
	}

	taskTypeValue, err := domain.NewTaskType(task.TaskType)
	if err != nil {
		return nil, err
	}

	return &domain.UserTask{
		ID:          taskID,
		UserID:      domainUserID,
		TaskType:    taskTypeValue,
		CompletedAt: task.CompletedAt,
		Points:      task.Points,
	}, nil
}
