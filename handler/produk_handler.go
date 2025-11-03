package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"rakamin-evermos/model"
	"rakamin-evermos/repository"
	"rakamin-evermos/usecase"
	"rakamin-evermos/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InputProduk struct {
	NamaProduk    string `json:"nama_produk" binding:"required"`
	Slug          string `json:"slug" binding:"required"`
	HargaReseller string `json:"harga_reseller" binding:"required"`
	HargaKonsumen string `json:"harga_konsumen" binding:"required"`
	Stok          int    `json:"stok" binding:"required"`
	Deskripsi     string `json:"deskripsi" binding:"required"`
	IDCategory    uint   `json:"id_category" binding:"required"`
}

type ProdukHandler interface {

	// Publik
	GetAllProduk(c *gin.Context)
	GetProdukByID(c *gin.Context)

	// Seller
	CreateProduk(c *gin.Context)
	GetMyProduk(c *gin.Context)
	UpdateProduk(c *gin.Context)
	DeleteProduk(c *gin.Context)
	UploadFotoProduk(c *gin.Context)
}

type produkHandler struct {
	produkUsecase usecase.ProdukUsecase
}

func NewProdukHandler(produkUsecase usecase.ProdukUsecase) ProdukHandler {
	return &produkHandler{produkUsecase}
}

//public accessible

func parseFilterAndPagination(c *gin.Context) (utils.PaginationInput, repository.FilterInput) {
	pagination := utils.GetPaginationFromQuery(c)

	// get filter
	search := c.Query("search")
	categoryID, _ := strconv.Atoi(c.Query("category_id"))

	filter := repository.FilterInput{
		Search:     search,
		CategoryID: uint(categoryID),
	}

	return pagination, filter
}

func (h *produkHandler) GetAllProduk(c *gin.Context) {
	pagination, filter := parseFilterAndPagination(c)

	result, err := h.produkUsecase.GetAllProduk(pagination, filter)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, "Success get all produk", result)
}

func (h *produkHandler) GetProdukByID(c *gin.Context) {
	produkID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID produk not valid")
		return
	}

	produk, err := h.produkUsecase.GetProdukByID(uint(produkID))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SendSuccessResponse(c, "Success get Detail produk", produk)
}

// seller only

func (h *produkHandler) CreateProduk(c *gin.Context) {
	userID, _ := c.Get("currentUserID")

	var input InputProduk
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	produk := model.Produk{
		NamaProduk:    input.NamaProduk,
		Slug:          input.Slug,
		HargaReseller: input.HargaReseller,
		HargaKonsumen: input.HargaKonsumen,
		Stok:          input.Stok,
		Deskripsi:     input.Deskripsi,
		IDCategory:    input.IDCategory,
	}

	savedProduk, err := h.produkUsecase.CreateProduk(userID.(uint), produk)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendCreatedResponse(c, "Success create produk", savedProduk)
}

func (h *produkHandler) GetMyProduk(c *gin.Context) {
	userID, _ := c.Get("currentUserID")

	pagination, filter := parseFilterAndPagination(c)

	result, err := h.produkUsecase.GetMyProduk(userID.(uint), pagination, filter)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, "Success get my produk", result)
}

func (h *produkHandler) UpdateProduk(c *gin.Context) {
	userID, _ := c.Get("currentUserID")
	produkID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID produk not valid")
		return
	}

	var input InputProduk
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	produk := model.Produk{
		NamaProduk:    input.NamaProduk,
		Slug:          input.Slug,
		HargaReseller: input.HargaReseller,
		HargaKonsumen: input.HargaKonsumen,
		Stok:          input.Stok,
		Deskripsi:     input.Deskripsi,
		IDCategory:    input.IDCategory,
	}

	updatedProduk, err := h.produkUsecase.UpdateProduk(userID.(uint), uint(produkID), produk)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SendSuccessResponse(c, "Success update produk", updatedProduk)
}

func (h *produkHandler) DeleteProduk(c *gin.Context) {
	userID, _ := c.Get("currentUserID")
	produkID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID produk not valid")
		return
	}

	err = h.produkUsecase.DeleteProduk(userID.(uint), uint(produkID))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SendSuccessResponse(c, "Success delete produk", nil)
}

func (h *produkHandler) UploadFotoProduk(c *gin.Context) {
	userID, _ := c.Get("currentUserID")
	produkID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID produk not valid")
		return
	}

	file, err := c.FormFile("photo")
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "File upload not found (key must be 'photo')")
		return
	}

	// create unique file uploads/produk-[produkID]-[uuid].[ext]
	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("produk-%d-%s%s", produkID, uuid.New().String(), ext)
	filePath := "uploads/" + fileName

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to save file")
		return
	}

	savedFoto, err := h.produkUsecase.UploadFotoProduk(userID.(uint), uint(produkID), filePath)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendCreatedResponse(c, "Success upload foto produk", savedFoto)
}