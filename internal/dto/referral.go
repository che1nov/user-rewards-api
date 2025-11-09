package dto

// ProcessReferralInput входные данные для обработки реферального кода
type ProcessReferralInput struct {
	ReferrerID string `json:"referrer_id" binding:"required"`
}

// ProcessReferralOutput выходные данные после обработки реферального кода
type ProcessReferralOutput struct {
	ReferralID     string `json:"referral_id"`
	ReferrerID     string `json:"referrer_id"`
	ReferredUserID string `json:"referred_user_id"`
	BonusPoints    int    `json:"bonus_points"`
	NewBalance     int    `json:"new_balance"`
}

