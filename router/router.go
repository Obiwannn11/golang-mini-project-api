package router

import (
	"net/http"

	"rakamin-evermos/handler"
	"rakamin-evermos/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, authHandler handler.AuthHandler, userHandler handler.UserHandler, addressHandler handler.AddressHandler) {

	api := r.Group("/api/v1")

	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login)

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

		authenticated.GET("users/me", userHandler.GetProfile)
		authenticated.PUT("users/me", userHandler.UpdateProfile)

		authenticated.POST("/addresses", addressHandler.CreateAddress)
		authenticated.GET("/addresses", addressHandler.GetAddresses)
		authenticated.GET("/addresses/:id", addressHandler.GetAddressByID)
		authenticated.PUT("/addresses/:id", addressHandler.UpdateAddress)
		authenticated.DELETE("/addresses/:id", addressHandler.DeleteAddress)

	}

	admin := api.Group("")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware())
	{

	}

}