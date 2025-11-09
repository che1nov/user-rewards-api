package dto

// LeaderboardEntry запись в таблице лидеров
type LeaderboardEntry struct {
	Rank     int    `json:"rank"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Balance  int    `json:"balance"`
}

// GetLeaderboardOutput выходные данные для таблицы лидеров
type GetLeaderboardOutput struct {
	Users []LeaderboardEntry `json:"users"`
	Total int                 `json:"total"`
}

