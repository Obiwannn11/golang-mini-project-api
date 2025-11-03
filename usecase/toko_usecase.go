package usecase

import (
	"errors"
	"fmt"
	"time"

	"rakamin-evermos/model"
	"rakamin-evermos/repository"

	"gorm.io/gorm"
)

type TokoUsecase interface {
	GetMyToko(userID uint) (model.Toko, error)
	UpdateMyToko(userID uint, input model.Toko) (model.Toko, error)
	UploadTokoPhoto(userID uint, filePath string) (model.Toko, error)
}

type tokoUsecase struct {
	tokoRepo repository.TokoRepository
}

func NewTokoUsecase(tokoRepo repository.TokoRepository) TokoUsecase {
	return &tokoUsecase{tokoRepo}
}


func (uc *tokoUsecase) GetMyToko(userID uint) (model.Toko, error) {
	toko, err := uc.tokoRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return toko, errors.New("toko not found")
		}
		return toko, fmt.Errorf("failed to retrieve toko data: %w", err)
	}
	return toko, nil
}

func (uc *tokoUsecase) UpdateMyToko(userID uint, input model.Toko) (model.Toko, error) {
	existingToko, err := uc.tokoRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Toko{}, errors.New("toko not found or you do not have access")
		}
		return model.Toko{}, fmt.Errorf("failed to verify toko: %w", err)
	}

	existingToko.NamaToko = input.NamaToko
	existingToko.UpdatedAtDate = time.Now()

	updatedToko, err := uc.tokoRepo.Update(existingToko)
	if err != nil {
		return updatedToko, fmt.Errorf("failed to update toko: %w", err)
	}
	return updatedToko, nil
}

func (uc *tokoUsecase) UploadTokoPhoto(userID uint, filePath string) (model.Toko, error) {
	existingToko, err := uc.tokoRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Toko{}, errors.New("toko not found or you do not have access")
		}
		return model.Toko{}, fmt.Errorf("failed to verify toko: %w", err)
	}

	existingToko.UrlFoto = filePath
	existingToko.UpdatedAtDate = time.Now()

	updatedToko, err := uc.tokoRepo.Update(existingToko)
	if err != nil {
		return updatedToko, fmt.Errorf("failed to upload toko photo: %w", err)
	}
	return updatedToko, nil
}