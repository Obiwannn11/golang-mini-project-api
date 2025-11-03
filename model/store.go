package model

import "time"

type Toko struct {
	ID            uint   `gorm:"primaryKey;autoIncrement;column:id"`
	IDUser        uint   `gorm:"column:id_user;unique"`
	NamaToko      string `gorm:"size:255"`
	UrlFoto       string `gorm:"size:255"`
	CreatedAtDate time.Time `gorm:"column:created_at_date"`
	UpdatedAtDate time.Time `gorm:"column:updated_at_date"`

	// Relasi ke produk
	Produk      []Produk `gorm:"foreignKey:IDToko"`
}

func (Toko) TableName() string {
	return "toko"
}