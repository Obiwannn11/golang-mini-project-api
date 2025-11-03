package usecase

import (
	"errors"
	"fmt"
	"time"

	"rakamin-evermos/model" 
	"rakamin-evermos/repository"
)

type UserUsecase interface {
	GetProfile(userID uint) (model.User, error)
	UpdateProfile(userID uint, updatedUser model.User) (model.User, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{userRepo}
}


func (uc *userUsecase) GetProfile(userID uint) (model.User, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return model.User{}, errors.New("profil user tidak ditemukan")
	}
	return user, nil
}

func (uc *userUsecase) UpdateProfile(userID uint, updatedUser model.User) (model.User, error) {
	existingUser, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return model.User{}, errors.New("profile user cant be found")
	}

	// the fields can be updated
	existingUser.Nama = updatedUser.Nama
	existingUser.NoTelp = updatedUser.NoTelp
	existingUser.TanggalLahir = updatedUser.TanggalLahir
	existingUser.JenisKelamin = updatedUser.JenisKelamin
	existingUser.Tentang = updatedUser.Tentang
	existingUser.Pekerjaan = updatedUser.Pekerjaan
	existingUser.IDProvinsi = updatedUser.IDProvinsi
	existingUser.IDKota = updatedUser.IDKota

    existingUser.UpdatedAtDate = time.Now()

	savedUser, err := uc.userRepo.Update(existingUser)
	if err != nil {
		return model.User{}, fmt.Errorf("failed update profile: %w", err)
	}

	return savedUser, nil
}