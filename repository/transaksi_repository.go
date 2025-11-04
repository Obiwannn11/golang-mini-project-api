package repository

import (
	"rakamin-evermos/model"

	"gorm.io/gorm"
)

type TransaksiRepository interface {
	Save(tx *gorm.DB, trx model.Trx) (model.Trx, error)

	// for see history
	FindAllByUserID(userID uint) ([]model.Trx, error)
	FindByUserAndTrxID(userID, trxID uint) (model.Trx, error)
}

type transaksiRepository struct {
	db *gorm.DB
}

func NewTransaksiRepository(db *gorm.DB) TransaksiRepository {
	return &transaksiRepository{db}
}


func (r *transaksiRepository) Save(tx *gorm.DB, trx model.Trx) (model.Trx, error) {
	err := tx.Create(&trx).Error
	return trx, err
}

func (r *transaksiRepository) FindAllByUserID(userID uint) ([]model.Trx, error) {
	var trxs []model.Trx
	// get all transaksi and also pre load detail
	err := r.db.Preload("DetailTrx").Preload("DetailTrx.LogProduk").Where("id_user = ?", userID).Find(&trxs).Error
	return trxs, err
}

func (r *transaksiRepository) FindByUserAndTrxID(userID, trxID uint) (model.Trx, error) {
	var trx model.Trx
	err := r.db.Preload("DetailTrx").Preload("DetailTrx.LogProduk").Where("id = ? AND id_user = ?", trxID, userID).First(&trx).Error
	return trx, err
}