package domain

type MutationsTable struct {
	ID        string  `gorm:"column:id;PRIMARY_KEY;NOT NULL"`
	UserID    string  `gorm:"column:user_id;NOT NULL"`
	Mutation  float64 `gorm:"column:mutation;NOT NULL"`
	Balance   float64 `gorm:"column:balance;NOT NULL"`
	Type      string  `gorm:"column:type;NOT NULL"`
	DepositID string  `gorm:"column:deposit_id"`
	ChargeID  string  `gorm:"column:charge_id"`
}

func (m *MutationsTable) TableName() string {
	return "mutations_table"
}
