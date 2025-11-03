package handler

import (
	"net/http"
	"strconv"

	"rakamin-evermos/model"
	"rakamin-evermos/usecase"
	"rakamin-evermos/utils"

	"github.com/gin-gonic/gin"
)

type CategoryInput struct {
	NamaCategory string `json:"nama_category" binding:"required"`
}

type CategoryHandler interface {
	CreateCategory(c *gin.Context)
	GetAllCategories(c *gin.Context)
	GetCategoryByID(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
}

type categoryHandler struct {
	categoryUsecase usecase.CategoryUsecase
}

func NewCategoryHandler(categoryUsecase usecase.CategoryUsecase) CategoryHandler {
	return &categoryHandler{categoryUsecase}
}

func (h *categoryHandler) CreateCategory(c *gin.Context) {
	var input CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	category := model.Category{
		NamaCategory: input.NamaCategory,
	}

	savedCategory, err := h.categoryUsecase.CreateCategory(category)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendCreatedResponse(c, "Success create Kategori", savedCategory)
}

func (h *categoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.categoryUsecase.GetAllCategories()
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, "Success get all kategori ", categories)
}

func (h *categoryHandler) GetCategoryByID(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID kategori not valid")
		return
	}

	category, err := h.categoryUsecase.GetCategoryByID(uint(categoryID))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SendSuccessResponse(c, "Success get Detail kategori", category)
}

func (h *categoryHandler) UpdateCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID kategori not valid")
		return
	}

	var input CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	inputCategory := model.Category{
		NamaCategory: input.NamaCategory,
	}

	updatedCategory, err := h.categoryUsecase.UpdateCategory(uint(categoryID), inputCategory)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SendSuccessResponse(c, "Success update kategori", updatedCategory)
}

func (h *categoryHandler) DeleteCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID kategori not valid")
		return
	}

	err = h.categoryUsecase.DeleteCategory(uint(categoryID))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SendSuccessResponse(c, "Success delete kategori", nil)
}