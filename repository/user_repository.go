package repository

import (
	"rakamin-evermos/model" 

	"gorm.io/gorm"
)

type UserRepository interface {
	Save(user model.User) (model.User, error)
	FindByEmail(email string) (model.User, error)
	FindByID(userID uint) (model.User, error)
	Update(user model.User) (model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Save(user model.User) (model.User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) FindByEmail(email string) (model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) FindByID(userID uint) (model.User, error) {
	var user model.User
	// Eager Loading Toko relation
	err := r.db.Preload("Toko").Where("id = ?", userID).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) Update(user model.User) (model.User, error) {
	err := r.db.Save(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}