package repository

import (
	"rakamin-evermos/model"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Save(category model.Category) (model.Category, error)
	FindAll() ([]model.Category, error)
	FindByID(categoryID uint) (model.Category, error)
	Update(category model.Category) (model.Category, error)
	Delete(category model.Category) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) Save(category model.Category) (model.Category, error) {
	err := r.db.Create(&category).Error
	if err != nil {
		return category, err
	}
	return category, nil
}

func (r *categoryRepository) FindAll() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Find(&categories).Error
	if err != nil {
		return categories, err
	}
	return categories, nil
}

func (r *categoryRepository) FindByID(categoryID uint) (model.Category, error) {
	var category model.Category
	err := r.db.Where("id = ?", categoryID).First(&category).Error
	if err != nil {
		return category, err
	}
	return category, nil
}

func (r *categoryRepository) Update(category model.Category) (model.Category, error) {
	err := r.db.Save(&category).Error
	if err != nil {
		return category, err
	}
	return category, nil
}

func (r *categoryRepository) Delete(category model.Category) error {
	err := r.db.Delete(&category).Error
	if err != nil {
		return err
	}
	return nil
}