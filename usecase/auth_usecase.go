package usecase

import (
	"errors"
	"fmt" 
	"time" 

	
	"rakamin-evermos/model"
	"rakamin-evermos/repository"
	"rakamin-evermos/utils"
)

type AuthUsecase interface {
	Register(user model.User) (model.User, error)
	Login(email, password string) (string, error) // return JWT token
}

type authUsecase struct {
	userRepo repository.UserRepository
	tokoRepo repository.TokoRepository 
}

func NewAuthUsecase(userRepo repository.UserRepository, tokoRepo repository.TokoRepository) AuthUsecase {
	return &authUsecase{userRepo, tokoRepo}
}

func (uc *authUsecase) Register(user model.User) (model.User, error) {
	_, err := uc.userRepo.FindByEmail(user.Email)
	if err == nil {
		return model.User{}, errors.New("email already registered")
	}

	// role default is user 
	user.IsAdmin = false

	hashedPassword, err := utils.HashPassword(user.KataSandi)
	if err != nil {
		return model.User{}, fmt.Errorf("failed encrypt pswrd: %w", err)
	}
	user.KataSandi = hashedPassword

    now := time.Now()
    user.CreatedAtDate = now
    user.UpdatedAtDate = now

	savedUser, err := uc.userRepo.Save(user)
	if err != nil {
		return model.User{}, fmt.Errorf("failed save user: %w", err)
	}

	// create toko after create user
	newToko := model.Toko{
		IDUser:        savedUser.ID,
		NamaToko:      fmt.Sprintf("%s's Toko", savedUser.Nama), // toko default name
        UrlFoto:       "",
        CreatedAtDate: now,
        UpdatedAtDate: now,
	}
	_, err = uc.tokoRepo.Save(newToko)
	if err != nil {
		fmt.Printf("failed make toko for user %d: %v\n", savedUser.ID, err)
	}

	return savedUser, nil
}

func (uc *authUsecase) Login(email, password string) (string, error) {
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("password or email incorrect")
	}

	if !utils.CheckPasswordHash(password, user.KataSandi) {
		return "", errors.New("password or email incorrect")
	}

	token, err := utils.GenerateToken(user.ID, user.IsAdmin)
	if err != nil {
		return "", fmt.Errorf("failed create token: %w", err)
	}

	return token, nil
}