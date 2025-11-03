package usecase

import (
	"errors"
	"fmt"
	"rakamin-evermos/model"
	"rakamin-evermos/repository"
	"rakamin-evermos/utils"
	"time"

	"gorm.io/gorm"
)

type ProdukUsecase interface {
	// public accessible 
	GetAllProduk(pagination utils.PaginationInput, filter repository.FilterInput) (utils.PaginationResult, error)
	GetProdukByID(produkID uint) (model.Produk, error)

	// seller only
	CreateProduk(userID uint, input model.Produk) (model.Produk, error)
	GetMyProduk(userID uint, pagination utils.PaginationInput, filter repository.FilterInput) (utils.PaginationResult, error)
	UpdateProduk(userID, produkID uint, input model.Produk) (model.Produk, error)
	DeleteProduk(userID, produkID uint) error
	UploadFotoProduk(userID, produkID uint, filePath string) (model.FotoProduk, error)
}

type produkUsecase struct {
	produkRepo     repository.ProdukRepository
	fotoProdukRepo repository.FotoProdukRepository
	tokoRepo       repository.TokoRepository 
}

func NewProdukUsecase(produkRepo repository.ProdukRepository, fotoProdukRepo repository.FotoProdukRepository, tokoRepo repository.TokoRepository) ProdukUsecase {
	return &produkUsecase{produkRepo, fotoProdukRepo, tokoRepo}
}


// Public Accessible

// get all produk with pagination & filtering
func (uc *produkUsecase) GetAllProduk(pagination utils.PaginationInput, filter repository.FilterInput) (utils.PaginationResult, error) {
	produks, totalData, err := uc.produkRepo.FindAll(pagination, filter)
	if err != nil {
		return utils.PaginationResult{}, fmt.Errorf("failed get produk: %w", err)
	}

	// format result from utils
	result := utils.GeneratePaginationResult(produks, totalData, pagination.Page, pagination.Limit)
	return result, nil
}

// get detail one produk
func (uc *produkUsecase) GetProdukByID(produkID uint) (model.Produk, error) {
	produk, err := uc.produkRepo.FindByID(produkID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return produk, errors.New("produk not found")
		}
		return produk, fmt.Errorf("failed get produk: %w", err)
	}
	return produk, nil
}

 // seller only

func (uc *produkUsecase) getTokoByUserID(userID uint) (model.Toko, error) {
	toko, err := uc.tokoRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return toko, errors.New("u dont have toko. go register as seller first")
		}
		return toko, fmt.Errorf("failed verify your toko: %w", err)
	}
	return toko, nil
}

func (uc *produkUsecase) CreateProduk(userID uint, input model.Produk) (model.Produk, error) {
	// get toko user first
	toko, err := uc.getTokoByUserID(userID)
	if err != nil {
		return model.Produk{}, err
	}

	// Set IDToko base toko login
	input.IDToko = toko.ID

	now := time.Now()
	input.CreatedAtDate = now
	input.UpdatedAtDate = now

	savedProduk, err := uc.produkRepo.Save(input)
	if err != nil {
		return savedProduk, fmt.Errorf("failed save produk: %w", err)
	}
	return savedProduk, nil
}

// get produk owned by toko user with pagination & filtering
func (uc *produkUsecase) GetMyProduk(userID uint, pagination utils.PaginationInput, filter repository.FilterInput) (utils.PaginationResult, error) {
	// get toko user first
	toko, err := uc.getTokoByUserID(userID)
	if err != nil {
		return utils.PaginationResult{}, err
	}

	produks, totalData, err := uc.produkRepo.FindAllByTokoID(toko.ID, pagination, filter)
	if err != nil {
		return utils.PaginationResult{}, fmt.Errorf("failed get your produk: %w", err)
	}

	// Format result
	result := utils.GeneratePaginationResult(produks, totalData, pagination.Page, pagination.Limit)
	return result, nil
}

func (uc *produkUsecase) UpdateProduk(userID, produkID uint, input model.Produk) (model.Produk, error) {
	toko, err := uc.getTokoByUserID(userID)
	if err != nil {
		return model.Produk{}, err
	}

	existingProduk, err := uc.produkRepo.FindByTokoIDAndProdukID(toko.ID, produkID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Produk{}, errors.New("produk not found atau u dont have access")
		}
		return model.Produk{}, fmt.Errorf("failed verify produk: %w", err)
	}

	existingProduk.NamaProduk = input.NamaProduk
	existingProduk.Slug = input.Slug
	existingProduk.HargaReseller = input.HargaReseller
	existingProduk.HargaKonsumen = input.HargaKonsumen
	existingProduk.Stok = input.Stok
	existingProduk.Deskripsi = input.Deskripsi
	existingProduk.IDCategory = input.IDCategory
	existingProduk.UpdatedAtDate = time.Now()

	updatedProduk, err := uc.produkRepo.Update(existingProduk)
	if err != nil {
		return updatedProduk, fmt.Errorf("failed update produk: %w", err)
	}
	return updatedProduk, nil
}

func (uc *produkUsecase) DeleteProduk(userID, produkID uint) error {
	toko, err := uc.getTokoByUserID(userID)
	if err != nil {
		return err
	}

	existingProduk, err := uc.produkRepo.FindByTokoIDAndProdukID(toko.ID, produkID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("produk not found atau u dont have access")
		}
		return fmt.Errorf("failed verify produk: %w", err)
	}

	if err := uc.produkRepo.Delete(existingProduk); err != nil {
		return fmt.Errorf("failed delete produk: %w", err)
	}
	return nil
}

func (uc *produkUsecase) UploadFotoProduk(userID, produkID uint, filePath string) (model.FotoProduk, error) {
	toko, err := uc.getTokoByUserID(userID)
	if err != nil {
		return model.FotoProduk{}, err
	}

	_, err = uc.produkRepo.FindByTokoIDAndProdukID(toko.ID, produkID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.FotoProduk{}, errors.New("produk not found atau u dont have access")
		}
		return model.FotoProduk{}, fmt.Errorf("failed verify produk: %w", err)
	}

	newFoto := model.FotoProduk{
		IDProduk:      produkID,
		Url:           filePath,
		CreatedAtDate: time.Now(),
		UpdatedAtDate: time.Now(),
	}

	savedFoto, err := uc.fotoProdukRepo.Save(newFoto)
	if err != nil {
		return savedFoto, fmt.Errorf("failed save foto produk: %w", err)
	}
	return savedFoto, nil
}