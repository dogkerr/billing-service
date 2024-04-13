package deposits

import "time"

type Deposit struct {
	ID        string    `gorm:"column:id;PRIMARY_KEY;NOT NULL"`
	UserID    string    `gorm:"column:user_id;NOT NULL"`
	Amount    float64   `gorm:"column:amount;NOT NULL"`
	Status    string    `gorm:"column:status;NOT NULL"`
	CreatedAt time.Time `gorm:"column:created_at;NOT NULL"`
	UpdatedAt time.Time `gorm:"column:updated_at;NOT NULL"`
}

func (m *Deposit) TableName() string {
	return "deposits_table"
}
