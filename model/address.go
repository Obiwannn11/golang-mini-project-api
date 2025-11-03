package model

import "time"

type Alamat struct {
	ID            uint   `gorm:"primaryKey;autoIncrement;column:id"`
	IDUser        uint   `gorm:"column:id_user"` 
	JudulAlamat   string `gorm:"size:255"`
	NamaPenerima  string `gorm:"size:255"`
	NoTelp        string `gorm:"size:255"`
	DetailAlamat  string `gorm:"size:255"`
	CreatedAtDate time.Time `gorm:"column:created_at_date"`
	UpdatedAtDate time.Time `gorm:"column:updated_at_date"`
}

func (Alamat) TableName() string {
	return "alamat"
}