package router

import (
	"net/http"

	"rakamin-evermos/handler"
	"rakamin-evermos/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine,
	 authHandler handler.AuthHandler,
	 userHandler handler.UserHandler,
	 addressHandler handler.AddressHandler,
	 categoryHandler handler.CategoryHandler,
	 tokoHandler handler.TokoHandler,
	 produkHandler handler.ProdukHandler,
	 transaksiHandler handler.TransaksiHandler,
) {

	api := r.Group("/api/v1")

	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login)

	api.GET("/produk", produkHandler.GetAllProduk)
	api.GET("/produk/:id", produkHandler.GetProdukByID)

	authenticated := api.Group("")
	authenticated.Use(middleware.AuthMiddleware())
	{
		// protected route example
		authenticated.GET("/test-auth", func(c *gin.Context) {
			userID, _ := c.Get("currentUserID")

			c.JSON(http.StatusOK, gin.H{
				"message": "You have accessed a protected route!",
				"user_id": userID,
			})
		})


		// User routes
		authenticated.GET("users/me", userHandler.GetProfile)
		authenticated.PUT("users/me", userHandler.UpdateProfile)

		// Address routes
		authenticated.POST("/addresses", addressHandler.CreateAddress)
		authenticated.GET("/addresses", addressHandler.GetAddresses)
		authenticated.GET("/addresses/:id", addressHandler.GetAddressByID)
		authenticated.PUT("/addresses/:id", addressHandler.UpdateAddress)
		authenticated.DELETE("/addresses/:id", addressHandler.DeleteAddress)

		// Toko routes
		authenticated.GET("/toko/me", tokoHandler.GetMyToko)
		authenticated.PUT("/toko/me", tokoHandler.UpdateMyToko)
		authenticated.POST("/toko/me/photo", tokoHandler.UploadTokoPhoto)

		// Produk routes
		authenticated.POST("/my-produk", produkHandler.CreateProduk)
		authenticated.GET("/my-produk", produkHandler.GetMyProduk)
		authenticated.PUT("/my-produk/:id", produkHandler.UpdateProduk)
		authenticated.DELETE("/my-produk/:id", produkHandler.DeleteProduk)
		authenticated.POST("/my-produk/:id/photo", produkHandler.UploadFotoProduk)

		// Transaksi routes
		authenticated.POST("/transaksi", transaksiHandler.CreateTransaksi) // Checkout
		authenticated.GET("/transaksi", transaksiHandler.GetMyTransaksi)   // history
		authenticated.GET("/transaksi/:id", transaksiHandler.GetMyTransaksiByID) // Detail history
	}

	admin := api.Group("")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware())
	{
		// Category routes
		admin.POST("/categories", categoryHandler.CreateCategory)
		admin.GET("/categories", categoryHandler.GetAllCategories)
		admin.GET("/categories/:id", categoryHandler.GetCategoryByID)
		admin.PUT("/categories/:id", categoryHandler.UpdateCategory)
		admin.DELETE("/categories/:id", categoryHandler.DeleteCategory)
	}

}