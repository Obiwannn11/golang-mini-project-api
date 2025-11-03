package model

import "time"

type Category struct {
	ID            uint   `gorm:"primaryKey;autoIncrement;column:id"`
	NamaCategory  string `gorm:"size:255"`
	CreatedAtDate time.Time `gorm:"column:created_at_date"`
	UpdatedAtDate time.Time `gorm:"column:updated_at_date"`

	// Relasi ke produk
	Produk      []Produk `gorm:"foreignKey:IDCategory"`
}

func (Category) TableName() string {
	return "category"
}