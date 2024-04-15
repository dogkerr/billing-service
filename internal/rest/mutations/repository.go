package mutations

import "gorm.io/gorm"

type Repository interface {
	GetMutations() ([]Mutation, error)
	GetMutationByID(id string) (Mutation, error)
	GetMutationByUserID(id string) ([]Mutation, error)
	CreateMutation(mutation Mutation) (Mutation, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetMutations() ([]Mutation, error) {
	var mutations []Mutation
	if err := r.db.Find(&mutations).Error; err != nil {
		return mutations, err
	}

	return mutations, nil
}

func (r *repository) GetMutationByUserID(id string) ([]Mutation, error) {
	var mutations []Mutation
	if err := r.db.Where("user_id = ?", id).Order("created_at desc").Find(&mutations).Error; err != nil {
		return mutations, err
	}

	return mutations, nil
}

func (r *repository) GetMutationByID(id string) (Mutation, error) {
	var mutation Mutation
	if err := r.db.Where("id = ?", id).First(&mutation).Error; err != nil {
		return mutation, err
	}

	return mutation, nil
}

func (r *repository) CreateMutation(mutation Mutation) (Mutation, error) {
	if err := r.db.Create(&mutation).Error; err != nil {
		return Mutation{}, err
	}

	return mutation, nil
}
