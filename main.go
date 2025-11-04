package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"rakamin-evermos/config"
	"rakamin-evermos/model"
	"rakamin-evermos/handler"
	"rakamin-evermos/repository"
	"rakamin-evermos/router"
	"rakamin-evermos/usecase"
)

var (
	db *gorm.DB = config.ConnectDB()
)

func main() {
	log.Println("Running Database Migration...")
	err := db.AutoMigrate(
		&model.User{},
		&model.Alamat{},
		&model.Toko{},
		&model.Category{},
		&model.Produk{},
		&model.FotoProduk{},
		&model.LogProduk{},
		&model.Trx{},
		&model.DetailTrx{},
	)
	if err != nil {
		log.Fatal("failed migrasi database:", err)
	}
	log.Println("Migrasi Database finished.")

	// gin router
	r := gin.Default()

	userRepo := repository.NewUserRepository(db)
	tokoRepo := repository.NewTokoRepository(db)
	addressRepo := repository.NewAddressRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	produkRepo := repository.NewProdukRepository(db)
	fotoProdukRepo := repository.NewFotoProdukRepository(db)
	transaksiRepo := repository.NewTransaksiRepository(db)
	detailTrxRepo := repository.NewDetailTrxRepository(db)
	logProdukRepo := repository.NewLogProdukRepository(db)

	authUsecase := usecase.NewAuthUsecase(userRepo, tokoRepo)
	userUsecase := usecase.NewUserUsecase(userRepo)
	addressUsecase := usecase.NewAddressUsecase(addressRepo)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	tokoUsecase := usecase.NewTokoUsecase(tokoRepo)
	produkUsecase := usecase.NewProdukUsecase(produkRepo, fotoProdukRepo, tokoRepo)
	transaksiUsecase := usecase.NewTransaksiUsecase(
		db,
		transaksiRepo,
		detailTrxRepo,
		logProdukRepo,
		produkRepo,
		addressRepo,
	)

	authHandler := handler.NewAuthHandler(authUsecase)
	userHandler := handler.NewUserHandler(userUsecase)
	addressHandler := handler.NewAddressHandler(addressUsecase)
	categoryHandler := handler.NewCategoryHandler(categoryUsecase)
	tokoHandler := handler.NewTokoHandler(tokoUsecase)
	produkHandler := handler.NewProdukHandler(produkUsecase)
	transaksiHandler := handler.NewTransaksiHandler(transaksiUsecase)

	router.SetupRouter(
		r,
		authHandler,
		userHandler,
		addressHandler,
		categoryHandler,
		tokoHandler,
		produkHandler,
		transaksiHandler,
)

	port := os.Getenv("PORT")
	log.Printf("Server running in http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed running server:", err)
	}
}