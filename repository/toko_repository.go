package repository

import (
	"rakamin-evermos/model" 

	"gorm.io/gorm"
)

type TokoRepository interface {
	Save(toko model.Toko) (model.Toko, error)
	FindByUserID(userID uint) (model.Toko, error)
	Update(toko model.Toko) (model.Toko, error)
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

func (r *tokoRepository) FindByUserID(userID uint) (model.Toko, error) {
	var toko model.Toko
	err := r.db.Where("id_user = ?", userID).First(&toko).Error
	if err != nil {
		return toko, err
	}
	return toko, nil
}

func (r *tokoRepository) Update(toko model.Toko) (model.Toko, error) {
	err := r.db.Save(&toko).Error
	if err != nil {
		return toko, err
	}
	return toko, nil
}