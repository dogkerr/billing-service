package charges

import "gorm.io/gorm"

type Repository interface {
	Save(charge Charge) (Charge, error)
	FindByID(id string) (Charge, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(charge Charge) (Charge, error) {
	if err := r.db.Create(&charge).Error; err != nil {
		return Charge{}, err
	}

	return charge, nil
}

func (r *repository) FindByID(id string) (Charge, error) {
	var charge Charge
	if err := r.db.Where("id = ?", id).First(&charge).Error; err != nil {
		return charge, err
	}

	return charge, nil
}
