package repository

import (
	"rakamin-evermos/model" 

	"gorm.io/gorm"
)

type TokoRepository interface {
	Save(toko model.Toko) (model.Toko, error)
}

type tokoRepository struct {
	db *gorm.DB
}

func NewTokoRepository(db *gorm.DB) TokoRepository {
	return &tokoRepository{db}
}

func (r *tokoRepository) Save(toko model.Toko) (model.Toko, error) {
	err := r.db.Create(&toko).Error
	if err != nil {
		return toko, err
	}
	return toko, nil
}