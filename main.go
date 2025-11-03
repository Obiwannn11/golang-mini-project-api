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

	authUsecase := usecase.NewAuthUsecase(userRepo, tokoRepo)
	userUsecase := usecase.NewUserUsecase(userRepo)

	authHandler := handler.NewAuthHandler(authUsecase)
	userHandler := handler.NewUserHandler(userUsecase)

	router.SetupRouter(r, authHandler, userHandler)

	port := os.Getenv("PORT")
	log.Printf("Server running in http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed running server:", err)
	}
}