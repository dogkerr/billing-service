package payment

import "time"

type Transaction struct {
	ID          string
	UserID      string
	Amount      int
	Token       string
	RedirectURL string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
