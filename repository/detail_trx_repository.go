package repository

import (
	"rakamin-evermos/model"

	"gorm.io/gorm"
)

type DetailTrxRepository interface {
	Save(tx *gorm.DB, detailTrx model.DetailTrx) (model.DetailTrx, error)
}

type detailTrxRepository struct {
	db *gorm.DB
}

func NewDetailTrxRepository(db *gorm.DB) DetailTrxRepository {
	return &detailTrxRepository{db}
}


func (r *detailTrxRepository) Save(tx *gorm.DB, detailTrx model.DetailTrx) (model.DetailTrx, error) {
	err := tx.Create(&detailTrx).Error
	return detailTrx, err
}