package repository

import (
	"rakamin-evermos/model"

	"gorm.io/gorm"
)

type LogProdukRepository interface {
	Save(tx *gorm.DB, logProduk model.LogProduk) (model.LogProduk, error)
}

type logProdukRepository struct {
	db *gorm.DB
}

func NewLogProdukRepository(db *gorm.DB) LogProdukRepository {
	return &logProdukRepository{db}
}


func (r *logProdukRepository) Save(tx *gorm.DB, logProduk model.LogProduk) (model.LogProduk, error) {
	// use tx from usecase not r db
	err := tx.Create(&logProduk).Error
	return logProduk, err
}