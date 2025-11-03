package model

import "time"

type FotoProduk struct {
	ID            uint   `gorm:"primaryKey;autoIncrement;column:id"`
	IDProduk      uint   `gorm:"column:id_produk"`
	Url           string `gorm:"size:255"`
	CreatedAtDate time.Time `gorm:"column:created_at_date"`
	UpdatedAtDate time.Time `gorm:"column:updated_at_date"`
}

func (FotoProduk) TableName() string {
	return "foto_produk"
}