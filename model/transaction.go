package model

import "time"

type Trx struct {
	ID               uint   `gorm:"primaryKey;autoIncrement;column:id"`
	IDUser           uint   `gorm:"column:id_user"`
	AlamatPengiriman uint   `gorm:"column:alamat_pengiriman"`
	HargaTotal       int
	KodeInvoice      string `gorm:"size:255"`
	MethodBayar      string `gorm:"size:255"`
	CreatedAtDate    time.Time `gorm:"column:created_at_date"`
	UpdatedAtDate    time.Time `gorm:"column:updated_at_date"`

	// Relasi yg ke detail transaksi, alamat
	DetailTrx      []DetailTrx `gorm:"foreignKey:IDTrx"`
	Alamat           Alamat      `gorm:"foreignKey:AlamatPengiriman"`
}

func (Trx) TableName() string {
	return "trx"
}