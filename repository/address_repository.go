package repository

import (
	"rakamin-evermos/model"

	"gorm.io/gorm"
)

type AddressRepository interface {
	Save(alamat model.Alamat) (model.Alamat, error)
	FindAllByUserID(userID uint) ([]model.Alamat, error)
	FindByIDAndUserID(addressID, userID uint) (model.Alamat, error)
	Update(alamat model.Alamat) (model.Alamat, error)
	Delete(alamat model.Alamat) error
}

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{db}
}


func (r *addressRepository) Save(alamat model.Alamat) (model.Alamat, error) {
	err := r.db.Create(&alamat).Error
	if err != nil {
		return alamat, err
	}
	return alamat, nil
}

// get all alamat based on ID User
func (r *addressRepository) FindAllByUserID(userID uint) ([]model.Alamat, error) {
	var alamats []model.Alamat
	err := r.db.Where("id_user = ?", userID).Find(&alamats).Error
	if err != nil {
		return alamats, err
	}
	return alamats, nil
}

// find alamat berdasarkan ID Alamat and ID User
func (r *addressRepository) FindByIDAndUserID(addressID, userID uint) (model.Alamat, error) {
	var alamat model.Alamat
	err := r.db.Where("id = ? AND id_user = ?", addressID, userID).First(&alamat).Error
	if err != nil {
		return alamat, err
	}
	return alamat, nil
}

func (r *addressRepository) Update(alamat model.Alamat) (model.Alamat, error) {
	err := r.db.Save(&alamat).Error
	if err != nil {
		return alamat, err
	}
	return alamat, nil
}

func (r *addressRepository) Delete(alamat model.Alamat) error {
	err := r.db.Delete(&alamat).Error
	if err != nil {
		return err
	}
	return nil
}