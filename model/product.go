package model

import "time"

type Produk struct {
	ID             uint   `gorm:"primaryKey;autoIncrement;column:id"`
	IDToko         uint   `gorm:"column:id_toko"`
	IDCategory     uint   `gorm:"column:id_category"`
	NamaProduk     string `gorm:"size:255"`
	Slug           string `gorm:"size:255"`
	HargaReseller  string `gorm:"size:255"`
	HargaKonsumen  string `gorm:"size:255"`
	Stok           int
	Deskripsi      string `gorm:"type:text"`
	CreatedAtDate  time.Time `gorm:"column:created_at_date"`
	UpdatedAtDate  time.Time `gorm:"column:updated_at_date"`

	// Relasi nya ke foto produk, log produk, dan kategori
	FotoProduk   []FotoProduk `gorm:"foreignKey:IDProduk"`
	LogProduk    []LogProduk  `gorm:"foreignKey:IDProduk"`
	Category     Category     `gorm:"foreignKey:IDCategory"`
}

func (Produk) TableName() string {
	return "produk"
}