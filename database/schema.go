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
	Id          string      `json:"id"`
	CreatedBy   string      `json:"created_by"`
	CreditsUsed int         `json:"credits_used"`
	Timestamp   time.Time   `json:"timestamp"`
	MediaUri    string      `json:"media_uri"`
	Type        string      `json:"type"`
	Meta        interface{} `json:"meta" db:"meta"`
}

type GenerationDB struct {
	Id          string    `json:"id"`
	CreatedBy   string    `json:"created_by"`
	CreditsUsed int       `json:"credits_used"`
	Timestamp   time.Time `json:"timestamp"`
	MediaUri    string    `json:"media_uri"`
	Type        string    `json:"type"`
	Meta        string    `json:"meta"`
}

type AiServices struct {
	Aid           string `json:"aid"`
	Type          string `json:"type"`
	Provider      string `json:"provider"`
	PrePrompt     string `json:"pre_prompt"`
	Banner        string `json:"banner"`
	Description   string `json:"description"`
	CreditsPerGen int    `json:"credits_per_gen"`
}

type Wallet struct {
	WalletId string `json:"wallet_id"`
	UserId   string `json:"user_id"`
	Balance  int    `json:"balance"`
}

type Transactions struct {
	TrxId       string    `json:"trx_id"`
	UserId      string    `json:"user_id"`
	WalletId    string    `json:"wallet_id"`
	Amount      int       `json:"amount"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

type Projects struct {
	ProjectId   string `json:"project_id"`
	Type        string `json:"type"`
	UserId      string `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Meta        string `json:"meta"`
	States      string `json:"states"`
}
