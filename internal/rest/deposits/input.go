package deposits

type DepositInput struct {
	ID     string  `json:"id" binding:"required"`
	UserID string  `json:"user_id" binding:"required"`
	Amount float32 `json:"amount" binding:"required"`
	Status string  `json:"status" binding:"required"`
}
