package deposits

type DepositInput struct {
	UserID string  `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}

type DepositUpdateInput struct {
	ID     string  `json:"id" binding:"required"`
	UserID string  `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
	Status string  `json:"status" binding:"required"`
}
