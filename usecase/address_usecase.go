package usecase

import (
	"errors"
	"fmt"
	"time"

	"rakamin-evermos/model"
	"rakamin-evermos/repository"

	"gorm.io/gorm"
)

type AddressUsecase interface {
	CreateAddress(userID uint, alamat model.Alamat) (model.Alamat, error)
	GetAddresses(userID uint) ([]model.Alamat, error)
	GetAddressByID(addressID, userID uint) (model.Alamat, error)
	UpdateAddress(addressID, userID uint, inputAlamat model.Alamat) (model.Alamat, error)
	DeleteAddress(addressID, userID uint) error
}

type addressUsecase struct {
	addressRepo repository.AddressRepository
}

func NewAddressUsecase(addressRepo repository.AddressRepository) AddressUsecase {
	return &addressUsecase{addressRepo}
}

func (uc *addressUsecase) CreateAddress(userID uint, alamat model.Alamat) (model.Alamat, error) {
	alamat.IDUser = userID

	now := time.Now()
	alamat.CreatedAtDate = now
	alamat.UpdatedAtDate = now

	savedAlamat, err := uc.addressRepo.Save(alamat)
	if err != nil {
		return savedAlamat, fmt.Errorf("failed save alamat: %w", err)
	}
	return savedAlamat, nil
}

func (uc *addressUsecase) GetAddresses(userID uint) ([]model.Alamat, error) {
	alamats, err := uc.addressRepo.FindAllByUserID(userID)
	if err != nil {
		return alamats, fmt.Errorf("failed get addresses: %w", err)
	}
	return alamats, nil
}

func (uc *addressUsecase) GetAddressByID(addressID, userID uint) (model.Alamat, error) {
	alamat, err := uc.addressRepo.FindByIDAndUserID(addressID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return alamat, errors.New("alamat tidak ditemukan atau anda tidak memiliki akses")
		}
		return alamat, fmt.Errorf("failed get address: %w", err)
	}
	return alamat, nil
}

func (uc *addressUsecase) UpdateAddress(addressID, userID uint, inputAlamat model.Alamat) (model.Alamat, error) {
	existingAlamat, err := uc.addressRepo.FindByIDAndUserID(addressID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Alamat{}, errors.New("alamat tidak ditemukan atau anda tidak memiliki akses")
		}
		return model.Alamat{}, fmt.Errorf("failed verify address: %w", err)
	}

	// field can be updated
	existingAlamat.JudulAlamat = inputAlamat.JudulAlamat
	existingAlamat.NamaPenerima = inputAlamat.NamaPenerima
	existingAlamat.NoTelp = inputAlamat.NoTelp
	existingAlamat.DetailAlamat = inputAlamat.DetailAlamat
	existingAlamat.UpdatedAtDate = time.Now()

	updatedAlamat, err := uc.addressRepo.Update(existingAlamat)
	if err != nil {
		return updatedAlamat, fmt.Errorf("failed update address: %w", err)
	}
	return updatedAlamat, nil
}

func (uc *addressUsecase) DeleteAddress(addressID, userID uint) error {
	existingAlamat, err := uc.addressRepo.FindByIDAndUserID(addressID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("alamat tidak ditemukan atau anda tidak memiliki akses")
		}
		return fmt.Errorf("failed verify address: %w", err)
	}

	err = uc.addressRepo.Delete(existingAlamat)
	if err != nil {
		return fmt.Errorf("failed delete address: %w", err)
	}
	return nil
}