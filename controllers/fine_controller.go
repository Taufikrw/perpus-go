package controllers

import (
	"belajar-go/dto"
	"belajar-go/resources"
	"belajar-go/services"
	"belajar-go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FineController struct {
	svc *services.FineService
}

func NewFineController(svc *services.FineService) *FineController {
	return &FineController{svc: svc}
}

func (ctrl *FineController) IndexFines(c *gin.Context) {
	fines, err := ctrl.svc.GetAllFines(c.Request.Context())
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Fines retrieved successfully!", resources.FormatFines(fines))
}

func (ctrl *FineController) ShowFine(c *gin.Context) {
	fine, err := ctrl.svc.GetFineByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Fine retrieved successfully!", resources.FormatFine(*fine))
}

func (ctrl *FineController) StoreFine(c *gin.Context) {
	var input dto.FineDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		utils.HandleError(c, utils.NewValidationError("Invalid input", errMsg))
		return
	}

	newFine, err := ctrl.svc.CreateFine(c.Request.Context(), input)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusCreated, "Fine created successfully!", resources.FormatFine(*newFine))
}

func (ctrl *FineController) UpdateFine(c *gin.Context) {
	var input dto.FineDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		utils.HandleError(c, utils.NewValidationError("Invalid input", errMsg))
		return
	}

	updatedFine, err := ctrl.svc.UpdateFine(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Fine updated successfully!", resources.FormatFine(*updatedFine))
}

func (ctrl *FineController) DeleteFine(c *gin.Context) {
	err := ctrl.svc.DeleteFine(c.Request.Context(), c.Param("id"))
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Fine deleted successfully!", nil)
}
