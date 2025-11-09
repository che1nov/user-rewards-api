package dto

// CreateUserInput входные данные для создания пользователя
type CreateUserInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

// CreateUserOutput выходные данные после создания пользователя
type CreateUserOutput struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

