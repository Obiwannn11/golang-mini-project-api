package handler

import (
	"net/http"
	"rakamin-evermos/usecase"
	"rakamin-evermos/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransaksiInput struct {
	AlamatPengirimanID uint                   `json:"alamat_pengiriman_id" binding:"required"`
	MethodBayar        string                 `json:"method_bayar" binding:"required"`
	Items              []usecase.CartItemInput `json:"items" binding:"required,dive"` // dive for vlidate nested array
}

type TransaksiHandler interface {
	CreateTransaksi(c *gin.Context)
	GetMyTransaksi(c *gin.Context)
	GetMyTransaksiByID(c *gin.Context)
}

type transaksiHandler struct {
	transaksiUsecase usecase.TransaksiUsecase
}

func NewTransaksiHandler(transaksiUsecase usecase.TransaksiUsecase) TransaksiHandler {
	return &transaksiHandler{transaksiUsecase}
}


// checkout / create transaksi
func (h *transaksiHandler) CreateTransaksi(c *gin.Context) {
	// 1. Ambil userID dari context
	userID, exists := c.Get("currentUserID")
	if !exists {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "fail get ID user from token")
		return
	}

	// validate input
	var input TransaksiInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	savedTrx, err := h.transaksiUsecase.CreateTransaksi(
		userID.(uint),
		input.AlamatPengirimanID,
		input.MethodBayar,
		input.Items,
	)
	if err != nil {
		// return 400, error is from invalid input ( invalid alamatID, produkID, stok not enough, etc)
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendCreatedResponse(c, "Success create Transaksi", savedTrx)
}

func (h *transaksiHandler) GetMyTransaksi(c *gin.Context) {
	userID, _ := c.Get("currentUserID")

	trxs, err := h.transaksiUsecase.GetMyTransaksi(userID.(uint))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, "Success get my transaksi", trxs)
}

func (h *transaksiHandler) GetMyTransaksiByID(c *gin.Context) {
	userID, _ := c.Get("currentUserID")

	// get trxID from URL
	trxID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID transaksi not valid")
		return
	}

	trx, err := h.transaksiUsecase.GetMyTransaksiByID(userID.(uint), uint(trxID))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SendSuccessResponse(c, "Success get Detail transaksi", trx)
}