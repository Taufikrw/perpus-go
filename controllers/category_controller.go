package controllers

import (
	"belajar-go/dto"
	"belajar-go/resources"
	"belajar-go/services"
	"belajar-go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	svc *services.CategoryService
}

func NewCategoryController(svc *services.CategoryService) *CategoryController {
	return &CategoryController{svc: svc}
}

func (ctrl *CategoryController) IndexCategory(c *gin.Context) {
	categories, err := ctrl.svc.GetAllCategories(c.Request.Context())
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Categories retrieved successfully!", resources.FormatCategories(categories))
}

func (ctrl *CategoryController) ShowCategory(c *gin.Context) {
	id := c.Param("id")

	category, err := ctrl.svc.GetCategoryByID(c.Request.Context(), id)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Category retrieved successfully!", resources.FormatCategory(*category))
}

func (ctrl *CategoryController) StoreCategory(c *gin.Context) {
	var input dto.CategoryDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		validationErr := utils.NewValidationError("Invalid Input", errMsg)

		utils.HandleError(c, validationErr)
		return
	}

	newCategory, err := ctrl.svc.CreateCategory(c.Request.Context(), input)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusCreated, "Category created successfully!", resources.FormatCategory(*newCategory))
}

func (ctrl *CategoryController) UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var input dto.CategoryDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		validationErr := utils.NewValidationError("Invalid Input", errMsg)
		utils.HandleError(c, validationErr)
		return
	}

	updatedCategory, err := ctrl.svc.UpdateCategory(c.Request.Context(), id, input)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Category updated successfully!", resources.FormatCategory(*updatedCategory))
}

func (ctrl *CategoryController) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	err := ctrl.svc.DeleteCategory(c.Request.Context(), id)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Category deleted successfully!", nil)
}
