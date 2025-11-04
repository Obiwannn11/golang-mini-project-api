package usecase

import (
	"errors"
	"fmt"
	"rakamin-evermos/model"
	"rakamin-evermos/repository"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartItemInput struct {
	ProdukID  uint `json:"produk_id"`
	Kuantitas int  `json:"kuantitas"`
}

type TransaksiUsecase interface {
	CreateTransaksi(userID, alamatID uint, methodBayar string, items []CartItemInput) (model.Trx, error)

	GetMyTransaksi(userID uint) ([]model.Trx, error)
	GetMyTransaksiByID(userID, trxID uint) (model.Trx, error)
}

type transaksiUsecase struct {
	db *gorm.DB

	transaksiRepo repository.TransaksiRepository
	detailTrxRepo repository.DetailTrxRepository
	logProdukRepo repository.LogProdukRepository
	produkRepo    repository.ProdukRepository
	addressRepo   repository.AddressRepository
}

func NewTransaksiUsecase(
	db *gorm.DB,
	transaksiRepo repository.TransaksiRepository,
	detailTrxRepo repository.DetailTrxRepository,
	logProdukRepo repository.LogProdukRepository,
	produkRepo repository.ProdukRepository,
	addressRepo repository.AddressRepository,
) TransaksiUsecase {
	return &transaksiUsecase{
		db,
		transaksiRepo,
		detailTrxRepo,
		logProdukRepo,
		produkRepo,
		addressRepo,
	}
}


func (uc *transaksiUsecase) CreateTransaksi(userID, alamatID uint, methodBayar string, items []CartItemInput) (model.Trx, error) {

	// verify userid and alamat
	_, err := uc.addressRepo.FindByIDAndUserID(alamatID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Trx{}, errors.New("alamat not found or access denied")
		}
		return model.Trx{}, err
	}

	tx := uc.db.Begin()
	if tx.Error != nil {
		return model.Trx{}, tx.Error
	}
	// defer for Rollback if panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var grandTotal int = 0
	var createdDetails []model.DetailTrx

	// Loop every item in cart
	for _, item := range items {
		// get produk and lock row
		produk, err := uc.produkRepo.FindByIDWithLock(tx, item.ProdukID)
		if err != nil {
			tx.Rollback()
			return model.Trx{}, errors.New("produk not found")
		}

		// check stok
		if produk.Stok < item.Kuantitas {
			tx.Rollback()
			return model.Trx{}, fmt.Errorf("stok for produk '%s' is not enough (remaining: %d)", produk.NamaProduk, produk.Stok)
		}

		// create log produk
		logProduk := model.LogProduk{
			IDProduk:      produk.ID,
			IDToko:        produk.IDToko,
			IDCategory:    produk.IDCategory,
			NamaProduk:    produk.NamaProduk,
			Slug:          produk.Slug,
			HargaReseller: produk.HargaReseller,
			HargaKonsumen: produk.HargaKonsumen,
			Deskripsi:     produk.Deskripsi,
			CreatedAtDate: time.Now(),
			UpdatedAtDate: time.Now(),
		}
		savedLog, err := uc.logProdukRepo.Save(tx, logProduk)
		if err != nil {
			tx.Rollback()
			return model.Trx{}, fmt.Errorf("fail save log produk: %w", err)
		}

		// count harga total per item and convert to varchar based on erd
		hargaItem, _ := strconv.Atoi(produk.HargaKonsumen)
		hargaTotalItem := hargaItem * item.Kuantitas
		grandTotal += hargaTotalItem

		// create detail trx
		detailTrx := model.DetailTrx{
			IDTrx:         0,
			IDLogProduk:   savedLog.ID,
			IDToko:        produk.IDToko,
			Kuantitas:     item.Kuantitas,
			HargaTotal:    hargaTotalItem,
			CreatedAtDate: time.Now(),
			UpdatedAtDate: time.Now(),
		}
		// save in array, then create after header
		createdDetails = append(createdDetails, detailTrx)

		// decrease Stok
		produk.Stok -= item.Kuantitas
		produk.UpdatedAtDate = time.Now()
		_, err = uc.produkRepo.UpdateWithTx(tx, produk)
		if err != nil {
			tx.Rollback()
			return model.Trx{}, fmt.Errorf("fail update stok: %w", err)
		}
	}

	// create Header Transaksi (Trx)
	newTrx := model.Trx{
		IDUser:           userID,
		AlamatPengiriman: alamatID,
		HargaTotal:       grandTotal,
		KodeInvoice:      fmt.Sprintf("INV/%d/%s", userID, uuid.New().String()[:8]), // make invoice unique code
		MethodBayar:      methodBayar,
		CreatedAtDate:    time.Now(),
		UpdatedAtDate:    time.Now(),
	}
	savedTrx, err := uc.transaksiRepo.Save(tx, newTrx)
	if err != nil {
		tx.Rollback()
		return model.Trx{}, fmt.Errorf("fail save header transaksi: %w", err)
	}

	// 5. Update IDTrx in all DetailTrx then save
	for _, detail := range createdDetails {
		detail.IDTrx = savedTrx.ID
		_, err := uc.detailTrxRepo.Save(tx, detail)
		if err != nil {
			tx.Rollback()
			return model.Trx{}, fmt.Errorf("fail save detail transaksi: %w", err)
		}
	}

	// Commit transaksi if success
	if err := tx.Commit().Error; err != nil {
		return model.Trx{}, fmt.Errorf("fail commit transaksi: %w", err)
	}

	// return header transaksi
	return savedTrx, nil
}


// history transaksi user
func (uc *transaksiUsecase) GetMyTransaksi(userID uint) ([]model.Trx, error) {
	trxs, err := uc.transaksiRepo.FindAllByUserID(userID)
	if err != nil {
		return trxs, fmt.Errorf("fail get history transaksi: %w", err)
	}
	return trxs, nil
}

// get detail transaksi by user
func (uc *transaksiUsecase) GetMyTransaksiByID(userID, trxID uint) (model.Trx, error) {
	trx, err := uc.transaksiRepo.FindByUserAndTrxID(userID, trxID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return trx, errors.New("transaksi not found or you don't have access")
		}
		return trx, fmt.Errorf("fail get detail transaksi: %w", err)
	}
	return trx, nil
}