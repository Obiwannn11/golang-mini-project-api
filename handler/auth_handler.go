package handler

import (
	"net/http"

	"rakamin-evermos/model"
	"rakamin-evermos/usecase"
	"rakamin-evermos/utils"

	"github.com/gin-gonic/gin"
)

// --- Input Structs (DTOs) ---
type RegisterInput struct {
	Nama         string `json:"nama" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	KataSandi    string `json:"kata_sandi" binding:"required,min=8"`
	NoTelp       string `json:"no_telp" binding:"required,min=10"`
}

type LoginInput struct {
	Email     string `json:"email" binding:"required,email"`
	KataSandi string `json:"kata_sandi" binding:"required"`
}

type AuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type authHandler struct {
	authUsecase usecase.AuthUsecase 
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) AuthHandler {
	return &authHandler{authUsecase}
}

func (h *authHandler) Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user := model.User{
		Nama:      input.Nama,
		Email:     input.Email,
		KataSandi: input.KataSandi,
		NoTelp:    input.NoTelp,
	}

	savedUser, err := h.authUsecase.Register(user)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendCreatedResponse(c, "Registrasi berhasil", savedUser)
}

func (h *authHandler) Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.authUsecase.Login(input.Email, input.KataSandi)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	data := gin.H{"token": token}
	utils.SendSuccessResponse(c, "Login berhasil", data)
}