package repository

import (
	"rakamin-evermos/model"

	"gorm.io/gorm"
)

// 1. Interface
type FotoProdukRepository interface {
	Save(fotoProduk model.FotoProduk) (model.FotoProduk, error)
}

type fotoProdukRepository struct {
	db *gorm.DB
}

func NewFotoProdukRepository(db *gorm.DB) FotoProdukRepository {
	return &fotoProdukRepository{db}
}


func (r *fotoProdukRepository) Save(fotoProduk model.FotoProduk) (model.FotoProduk, error) {
	err := r.db.Create(&fotoProduk).Error
	return fotoProduk, err
}