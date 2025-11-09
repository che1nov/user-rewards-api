package dto

// CompleteTaskInput входные данные для выполнения задания
type CompleteTaskInput struct {
	TaskType string `json:"task_type" binding:"required"`
}

// CompleteTaskOutput выходные данные после выполнения задания
type CompleteTaskOutput struct {
	TaskID     string `json:"task_id"`
	TaskType   string `json:"task_type"`
	Points     int    `json:"points"`
	NewBalance int    `json:"new_balance"`
}

