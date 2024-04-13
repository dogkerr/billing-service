package deposits

import "gorm.io/gorm"

type Repository interface {
	Save(deposit Deposit) (Deposit, error)
	FindByID(id string) (Deposit, error)
	Update(deposit Deposit) (Deposit, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(deposit Deposit) (Deposit, error) {
	if err := r.db.Create(&deposit).Error; err != nil {
		return Deposit{}, err
	}

	return deposit, nil
}

func (r *repository) FindByID(id string) (Deposit, error) {
	var deposit Deposit
	if err := r.db.Where("id = ?", id).First(&deposit).Error; err != nil {
		return deposit, err
	}

	return deposit, nil
}

func (r *repository) Update(deposit Deposit) (Deposit, error) {
	if err := r.db.Save(&deposit).Error; err != nil {
		return Deposit{}, err
	}

	return deposit, nil
}
