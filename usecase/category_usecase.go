package usecase

import (
	"errors"
	"fmt"
	"time"

	"rakamin-evermos/model"
	"rakamin-evermos/repository"

	"gorm.io/gorm"
)

type CategoryUsecase interface {
	CreateCategory(input model.Category) (model.Category, error)
	GetAllCategories() ([]model.Category, error)
	GetCategoryByID(categoryID uint) (model.Category, error)
	UpdateCategory(categoryID uint, input model.Category) (model.Category, error)
	DeleteCategory(categoryID uint) error
}

type categoryUsecase struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryUsecase(categoryRepo repository.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{categoryRepo}
}


func (uc *categoryUsecase) CreateCategory(input model.Category) (model.Category, error) {
	now := time.Now()
	input.CreatedAtDate = now
	input.UpdatedAtDate = now

	savedCategory, err := uc.categoryRepo.Save(input)
	if err != nil {
		return savedCategory, fmt.Errorf("failed save kategori: %w", err)
	}
	return savedCategory, nil
}

func (uc *categoryUsecase) GetAllCategories() ([]model.Category, error) {
	categories, err := uc.categoryRepo.FindAll()
	if err != nil {
		return categories, fmt.Errorf("failed get all kategori: %w", err)
	}
	return categories, nil
}

func (uc *categoryUsecase) GetCategoryByID(categoryID uint) (model.Category, error) {
	category, err := uc.categoryRepo.FindByID(categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return category, errors.New("kategori not found")
		}
		return category, fmt.Errorf("failed get kategori: %w", err)
	}
	return category, nil
}

func (uc *categoryUsecase) UpdateCategory(categoryID uint, input model.Category) (model.Category, error) {
	existingCategory, err := uc.categoryRepo.FindByID(categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Category{}, errors.New("kategori not found")
		}
		return model.Category{}, fmt.Errorf("failed verify kategori: %w", err)
	}

	existingCategory.NamaCategory = input.NamaCategory
	existingCategory.UpdatedAtDate = time.Now()

	updatedCategory, err := uc.categoryRepo.Update(existingCategory)
	if err != nil {
		return updatedCategory, fmt.Errorf("failed update kategori: %w", err)
	}
	return updatedCategory, nil
}

func (uc *categoryUsecase) DeleteCategory(categoryID uint) error {
	existingCategory, err := uc.categoryRepo.FindByID(categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("kategori not found")
		}
		return fmt.Errorf("failed verify kategori: %w", err)
	}

	err = uc.categoryRepo.Delete(existingCategory)
	if err != nil {
		return fmt.Errorf("failed delete kategori: %w", err)
	}
	return nil
}