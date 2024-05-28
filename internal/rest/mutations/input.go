package mutations

type MutationInput struct {
	ID        string  `json:"id" binding:"required"`
	UserID    string  `json:"user_id" binding:"required"`
	Mutation  float32 `json:"mutation" binding:"required"`
	Type      string  `json:"type" binding:"required"`
	DepositID string  `json:"deposit_id"`
	ChargeID  string  `json:"charge_id"`
}
