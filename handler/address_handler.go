package handler

import (
	"net/http"
	"strconv"

	"rakamin-evermos/model"
	"rakamin-evermos/usecase"
	"rakamin-evermos/utils"

	"github.com/gin-gonic/gin"
)

// define data for Create dan Update
type AddressInput struct {
	JudulAlamat  string `json:"judul_alamat" binding:"required"`
	NamaPenerima string `json:"nama_penerima" binding:"required"`
	NoTelp       string `json:"no_telp" binding:"required"`
	DetailAlamat string `json:"detail_alamat" binding:"required"`
}

type AddressHandler interface {
	CreateAddress(c *gin.Context)
	GetAddresses(c *gin.Context)
	GetAddressByID(c *gin.Context)
	UpdateAddress(c *gin.Context)
	DeleteAddress(c *gin.Context)
}

type addressHandler struct {
	addressUsecase usecase.AddressUsecase
}

func NewAddressHandler(addressUsecase usecase.AddressUsecase) AddressHandler {
	return &addressHandler{addressUsecase}
}


func (h *addressHandler) CreateAddress(c *gin.Context) {
	userID, _ := c.Get("currentUserID")

	var input AddressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// change DTO to model.Alamat
	alamat := model.Alamat{
		JudulAlamat:  input.JudulAlamat,
		NamaPenerima: input.NamaPenerima,
		NoTelp:       input.NoTelp,
		DetailAlamat: input.DetailAlamat,
	}

	savedAlamat, err := h.addressUsecase.CreateAddress(userID.(uint), alamat)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendCreatedResponse(c, "Success added Alamat", savedAlamat)
}

func (h *addressHandler) GetAddresses(c *gin.Context) {
	userID, _ := c.Get("currentUserID")

	alamats, err := h.addressUsecase.GetAddresses(userID.(uint))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, "Success get all alamat", alamats)
}

func (h *addressHandler) GetAddressByID(c *gin.Context) {
	userID, _ := c.Get("currentUserID")

	addressID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID alamat not valid")
		return
	}

	alamat, err := h.addressUsecase.GetAddressByID(uint(addressID), userID.(uint))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SendSuccessResponse(c, "Success get Alamat detail", alamat)
}

func (h *addressHandler) UpdateAddress(c *gin.Context) {
	// 1. Ambil userID dan addressID
	userID, _ := c.Get("currentUserID")
	addressID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID alamat not valid")
		return
	}

	var input AddressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Mapping DTO to Model
	inputAlamat := model.Alamat{
		JudulAlamat:  input.JudulAlamat,
		NamaPenerima: input.NamaPenerima,
		NoTelp:       input.NoTelp,
		DetailAlamat: input.DetailAlamat,
	}

	updatedAlamat, err := h.addressUsecase.UpdateAddress(uint(addressID), userID.(uint), inputAlamat)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SendSuccessResponse(c, "Success update Alamat", updatedAlamat)
}

func (h *addressHandler) DeleteAddress(c *gin.Context) {
	userID, _ := c.Get("currentUserID")
	addressID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID alamat not valid")
		return
	}

	err = h.addressUsecase.DeleteAddress(uint(addressID), userID.(uint))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SendSuccessResponse(c, "Success delete Alamat", nil)
}