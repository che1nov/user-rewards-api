package dto

// GetUserStatusOutput выходные данные для получения статуса пользователя
type GetUserStatusOutput struct {
	UserID         string `json:"user_id"`
	Balance        int    `json:"balance"`
	CompletedTasks int    `json:"completed_tasks"`
	ReferralCount  int    `json:"referral_count"`
}

