package handler

import (
	"fmt"
	"net/http"
	"path/filepath" 

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"rakamin-evermos/model"
	"rakamin-evermos/usecase"
	"rakamin-evermos/utils"
)

// define data can change (only name toko so far)
type UpdateTokoInput struct {
	NamaToko string `json:"nama_toko" binding:"required"`
}

type TokoHandler interface {
	GetMyToko(c *gin.Context)
	UpdateMyToko(c *gin.Context)
	UploadTokoPhoto(c *gin.Context)
}

type tokoHandler struct {
	tokoUsecase usecase.TokoUsecase 
}

func NewTokoHandler(tokoUsecase usecase.TokoUsecase) TokoHandler {
	return &tokoHandler{tokoUsecase}
}


func (h *tokoHandler) GetMyToko(c *gin.Context) {
	userID, exists := c.Get("currentUserID")
	if !exists {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "fail to get user ID from token")
		return
	}

	toko, err := h.tokoUsecase.GetMyToko(userID.(uint))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SendSuccessResponse(c, "Success get Data toko", toko)
}

func (h *tokoHandler) UpdateMyToko(c *gin.Context) {
	userID, exists := c.Get("currentUserID")
	if !exists {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "failed to get user ID from token")
		return
	}

	var input UpdateTokoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedToko := model.Toko{
		NamaToko: input.NamaToko,
	}

	savedToko, err := h.tokoUsecase.UpdateMyToko(userID.(uint), updatedToko)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, "Success update Data toko", savedToko)
}

func (h *tokoHandler) UploadTokoPhoto(c *gin.Context) {
	userID, exists := c.Get("currentUserID")
	if !exists {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "failed to get user ID from token")
		return
	}

	file, err := c.FormFile("photo")
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "File upload not found (key must be 'photo')")
		return
	}

	// create unique file name
	// Format: uploads/toko-[userID]-[uuid].[ext]
	// Contoh: uploads/toko-1-a8a5b2e5-b1a5-4e78-a83d-6b5c7e1b5b5a.jpg
	ext := filepath.Ext(file.Filename) 
	fileName := fmt.Sprintf("toko-%d-%s%s", userID.(uint), uuid.New().String(), ext)
	filePath := "uploads/" + fileName

	// save to uploads/ folder
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to save file")
		return
	}

	// Call Usecase with the file Complete Folder LOCATION
	updatedToko, err := h.tokoUsecase.UploadTokoPhoto(userID.(uint), filePath)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, "Success Upload Photo toko", updatedToko)
}

