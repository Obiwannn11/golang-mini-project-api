package router

import (
	"rakamin-evermos/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, authHandler handler.AuthHandler) {

	api := r.Group("/api/v1")

	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login)

}