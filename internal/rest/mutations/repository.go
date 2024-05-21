package mutations

import "gorm.io/gorm"

type Repository interface {
	FindByID(id string) (Mutation, error)
	FindByUserID(id string) ([]Mutation, error)
	Save(mutation Mutation) (Mutation, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindByID(id string) (Mutation, error) {
	var mutation Mutation
	if err := r.db.Where("id = ?", id).First(&mutation).Error; err != nil {
		return mutation, err
	}

	return mutation, nil
}

func (r *repository) FindByUserID(id string) ([]Mutation, error) {
	var mutations []Mutation
	if err := r.db.Where("user_id = ?", id).Order("created_at desc").Find(&mutations).Error; err != nil {
		return mutations, err
	}

	return mutations, nil
}

func (r *repository) Save(mutation Mutation) (Mutation, error) {
	if err := r.db.Create(&mutation).Error; err != nil {
		return Mutation{}, err
	}

	return mutation, nil
}
