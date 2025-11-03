package repository

import (
	"rakamin-evermos/model"
	"rakamin-evermos/utils"

	"gorm.io/gorm"
)

//  produck parameter filter 
type FilterInput struct {
	Search     string
	CategoryID uint
}

type ProdukRepository interface {
	Save(produk model.Produk) (model.Produk, error)
	Update(produk model.Produk) (model.Produk, error)
	Delete(produk model.Produk) error
	FindByID(produkID uint) (model.Produk, error)

	FindByTokoIDAndProdukID(tokoID, produkID uint) (model.Produk, error)

	FindAll(pagination utils.PaginationInput, filter FilterInput) ([]model.Produk, int64, error)
	FindAllByTokoID(tokoID uint, pagination utils.PaginationInput, filter FilterInput) ([]model.Produk, int64, error)
}

type produkRepository struct {
	db *gorm.DB
}

func NewProdukRepository(db *gorm.DB) ProdukRepository {
	return &produkRepository{db}
}

func (r *produkRepository) Save(produk model.Produk) (model.Produk, error) {
	err := r.db.Create(&produk).Error
	return produk, err
}

func (r *produkRepository) Update(produk model.Produk) (model.Produk, error) {
	err := r.db.Save(&produk).Error
	return produk, err
}

func (r *produkRepository) Delete(produk model.Produk) error {
	return r.db.Delete(&produk).Error
}

func (r *produkRepository) FindByID(produkID uint) (model.Produk, error) {
	var produk model.Produk
	// Preload Kategori and Toko for more data
	err := r.db.Preload("Category").Preload("Toko").Where("id = ?", produkID).First(&produk).Error
	return produk, err
}

func (r *produkRepository) FindByTokoIDAndProdukID(tokoID, produkID uint) (model.Produk, error) {
	var produk model.Produk
	err := r.db.Where("id = ? AND id_toko = ?", produkID, tokoID).First(&produk).Error
	return produk, err
}

// use filter to query GORM
func buildFilterQuery(db *gorm.DB, filter FilterInput) *gorm.DB {
	query := db
	if filter.Search != "" {
		query = query.Where("nama_produk LIKE ?", "%"+filter.Search+"%")
	}
	if filter.CategoryID != 0 {
		query = query.Where("id_category = ?", filter.CategoryID)
	}
	return query
}

func (r *produkRepository) FindAll(pagination utils.PaginationInput, filter FilterInput) ([]model.Produk, int64, error) {
	var produks []model.Produk
	var totalData int64

	// base query 
	query := r.db.Model(&model.Produk{})

	// apply Filter
	query = buildFilterQuery(query, filter)

	// count total data bfore pagination
	err := query.Count(&totalData).Error
	if err != nil {
		return produks, totalData, err
	}

	// apply Pagination from utils
	err = query.Scopes(utils.Paginate(pagination.Page, pagination.Limit)).Preload("Category").Preload("Toko").Find(&produks).Error

	return produks, totalData, err
}

func (r *produkRepository) FindAllByTokoID(tokoID uint, pagination utils.PaginationInput, filter FilterInput) ([]model.Produk, int64, error) {
	var produks []model.Produk
	var totalData int64

	query := r.db.Model(&model.Produk{}).Where("id_toko = ?", tokoID)

	query = buildFilterQuery(query, filter)

	err := query.Count(&totalData).Error
	if err != nil {
		return produks, totalData, err
	}

	err = query.Scopes(utils.Paginate(pagination.Page, pagination.Limit)).Preload("Category").Find(&produks).Error

	return produks, totalData, err
}