package model

import "time"

// LogProduk mewakili tabel 'log_produk'
type LogProduk struct {
	ID             uint   `gorm:"primaryKey;autoIncrement;column:id"`
	IDProduk       uint   `gorm:"column:id_produk"`
	IDToko         uint   `gorm:"column:id_toko"`
	IDCategory     uint   `gorm:"column:id_category"`
	NamaProduk     string `gorm:"size:255"`
	Slug           string `gorm:"size:255"`
	HargaReseller  string `gorm:"size:255"`
	HargaKonsumen  string `gorm:"size:255"`
	Deskripsi      string `gorm:"type:text"`
	CreatedAtDate  time.Time `gorm:"column:created_at_date"`
	UpdatedAtDate  time.Time `gorm:"column:updated_at_date"`

	// Relasi yg ke detail transaksi
	DetailTrx []DetailTrx `gorm:"foreignKey:IDLogProduk"`
}

func (LogProduk) TableName() string {
	return "log_produk"
}