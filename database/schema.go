package database

import "time"

type ResponseHTTP struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type Users struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Location     string `json:"location"`
	ReferralCode string `json:"referral_code"`
}

type Meta struct {
	ModelId        string `json:"model_id"`
	Dimensions     string `json:"dimensions"`
	Prompt         string `json:"prompt"`
	NegativePrompt string `json:"negative_prompt"`
	BatchSize      int    `json:"batch_size"`
}

type Generation struct {
	Id          string    `json:"id"`
	CreatedBy   string    `json:"created_by"`
	CreditsUsed int       `json:"credits_used"`
	Timestamp   time.Time `json:"timestamp"`
	MediaUri    string    `json:"media_uri"`
	Type        string    `json:"type"`
	Meta        Meta      `json:"meta"`
}
