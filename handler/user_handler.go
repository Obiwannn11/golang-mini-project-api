package handler

import (
	"net/http"
	"time"

	"rakamin-evermos/model" 
	"rakamin-evermos/usecase"
	"rakamin-evermos/utils" 

	"github.com/gin-gonic/gin"
)

// define data (DTO) can updated by user
type UpdateProfileInput struct {
	Nama         string     `json:"nama" binding:"required"`
	NoTelp       string     `json:"no_telp" binding:"required"`
	TanggalLahir *time.Time `json:"tanggal_lahir"` 
	JenisKelamin string     `json:"jenis_kelamin"`
	Tentang      string     `json:"tentang"`
	Pekerjaan    string     `json:"pekerjaan"`
	IDProvinsi   int        `json:"id_provinsi"`
	IDKota       int        `json:"id_kota"`
}

// define user response data to send back
type UserProfileResponse struct {
	ID           uint       `json:"id"`
	Nama         string     `json:"nama"`
	Email        string     `json:"email"`
	NoTelp       string     `json:"no_telp"`
	TanggalLahir *time.Time `json:"tanggal_lahir"`
	JenisKelamin string     `json:"jenis_kelamin"`
	Tentang      string     `json:"tentang"`
	Pekerjaan    string     `json:"pekerjaan"`
	IDProvinsi   int        `json:"id_provinsi"`
	IDKota       int        `json:"id_kota"`
	IsAdmin      bool       `json:"is_admin"`
	Toko         model.Toko `json:"toko"` // include toko data
}

type UserHandler interface {
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
}

type userHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) UserHandler {
	return &userHandler{userUsecase}
}

func (h *userHandler) GetProfile(c *gin.Context) {
	// take userID from context (inside middleware)
	userID, exists := c.Get("currentUserID")
	if !exists {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "failed get ID user from token")
		return
	}

	user, err := h.userUsecase.GetProfile(userID.(uint))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	//  response DTO (no password)
	response := UserProfileResponse{
		ID:           user.ID,
		Nama:         user.Nama,
		Email:        user.Email,
		NoTelp:       user.NoTelp,
		TanggalLahir: user.TanggalLahir,
		JenisKelamin: user.JenisKelamin,
		Tentang:      user.Tentang,
		Pekerjaan:    user.Pekerjaan,
		IDProvinsi:   user.IDProvinsi,
		IDKota:       user.IDKota,
		IsAdmin:      user.IsAdmin,
		Toko:         user.Toko,
	}

	utils.SendSuccessResponse(c, "Profil user berhasil didapatkan", response)
}

func (h *userHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("currentUserID")
	if !exists {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "failed get ID user from token")
		return
	}

	// valiudate input JSON
	var input UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// 3. change input DTO to model.User
	updatedUser := model.User{
		Nama:         input.Nama,
		NoTelp:       input.NoTelp,
		TanggalLahir: input.TanggalLahir,
		JenisKelamin: input.JenisKelamin,
		Tentang:      input.Tentang,
		Pekerjaan:    input.Pekerjaan,
		IDProvinsi:   input.IDProvinsi,
		IDKota:       input.IDKota,
	}

	savedUser, err := h.userUsecase.UpdateProfile(userID.(uint), updatedUser)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Format response DTO (no password)
	response := UserProfileResponse{
		ID:           savedUser.ID,
		Nama:         savedUser.Nama,
		Email:        savedUser.Email,
		NoTelp:       savedUser.NoTelp,
		TanggalLahir: savedUser.TanggalLahir,
		JenisKelamin: savedUser.JenisKelamin,
		Tentang:      savedUser.Tentang,
		Pekerjaan:    savedUser.Pekerjaan,
		IDProvinsi:   savedUser.IDProvinsi,
		IDKota:       savedUser.IDKota,
		IsAdmin:      savedUser.IsAdmin,
		Toko:         savedUser.Toko,
	}

	utils.SendSuccessResponse(c, "Profil user berhasil diperbarui", response)
}