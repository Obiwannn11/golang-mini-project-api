package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"rakamin-evermos/config"
	"rakamin-evermos/model"
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
		log.Fatal("Gagal melakukan migrasi database:", err)
	}
	log.Println("Migrasi Database Selesai.")

	r := gin.Default()

    // tes
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })

	port := os.Getenv("PORT")
	log.Printf("Server berjalan di http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}