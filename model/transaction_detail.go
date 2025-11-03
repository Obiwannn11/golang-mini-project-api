package model

import "time"

type DetailTrx struct {
	ID            uint `gorm:"primaryKey;autoIncrement;column:id"`
	IDTrx         uint `gorm:"column:id_trx"`
	IDLogProduk   uint `gorm:"column:id_log_produk"`
	IDToko        uint `gorm:"column:id_toko"`
	Kuantitas     int
	HargaTotal    int
	CreatedAtDate time.Time `gorm:"column:created_at_date"`
	UpdatedAtDate time.Time `gorm:"column:updated_at_date"`

	// Relasi ke transaksi, log produk, toko
	LogProduk LogProduk `gorm:"foreignKey:IDLogProduk"`
	Toko      Toko      `gorm:"foreignKey:IDToko"`
}

func (DetailTrx) TableName() string {
	return "detail_trx"
}