package model

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement;column:id"`
	Nama         string    `gorm:"size:255"`
	KataSandi    string    `gorm:"size:255"`
	NoTelp       string    `gorm:"size:255;unique"`
	TanggalLahir *time.Time `gorm:"type:date"`
	JenisKelamin string    `gorm:"size:255"`
	Tentang      string    `gorm:"type:text"`
	Pekerjaan    string    `gorm:"size:255"`
	Email        string    `gorm:"size:255;unique"`
	IDProvinsi   int       `gorm:"column:id_provinsi"`
	IDKota       int       `gorm:"column:id_kota"`
	IsAdmin      bool      `gorm:"default:false"`
	CreatedAtDate time.Time `gorm:"column:created_at_date"`
	UpdatedAtDate time.Time `gorm:"column:updated_at_date"`

	// Relasi nya ke toko, alamat, dan transaksi
	Toko    Toko    `gorm:"foreignKey:IDUser"` 
	Alamat  []Alamat `gorm:"foreignKey:IDUser"`
	Trx     []Trx    `gorm:"foreignKey:IDUser"`
}

// nama tabel 
func (User) TableName() string {
	return "user" // Sesuai ERD
}