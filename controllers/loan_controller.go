package controllers

import (
	"belajar-go/dto"
	"belajar-go/resources"
	"belajar-go/services"
	"belajar-go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoanController struct {
	service *services.LoanService
}

func NewLoanController(service *services.LoanService) *LoanController {
	return &LoanController{service: service}
}

func (ctrl *LoanController) IndexLoans(c *gin.Context) {
	loans, err := ctrl.service.GetAllLoans(c.Request.Context())
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Daftar peminjaman berhasil diambil!", resources.FormatLoans(loans))
}

func (ctrl *LoanController) StoreLoan(c *gin.Context) {
	var input dto.CreateLoanDTO
	userID, _ := c.Get("user_id")

	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		utils.HandleError(c, utils.NewValidationError("Invalid input", errMsg))
		return
	}

	newLoan, err := ctrl.service.CreateLoan(c.Request.Context(), input, userID.(string))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusCreated, "Loan created successfully!", resources.FormatLoan(*newLoan))
}

func (ctrl *LoanController) ShowLoan(c *gin.Context) {
	loan, err := ctrl.service.GetLoanByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Loan detail retrieved successfully!", resources.FormatLoan(*loan))
}

func (ctrl *LoanController) UpdateLoan(c *gin.Context) {
	var input dto.UpdateLoanDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		utils.HandleError(c, utils.NewValidationError("Invalid input", errMsg))
		return
	}

	updatedLoan, err := ctrl.service.UpdateLoan(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Loan updated successfully!", resources.FormatLoan(*updatedLoan))
}

func (ctrl *LoanController) ReturnLoan(c *gin.Context) {
	id := c.Param("id")
	loan, err := ctrl.service.ReturnLoan(c.Request.Context(), id)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Loan returned successfully!", resources.FormatLoan(*loan))
}

func (ctrl *LoanController) DeleteLoan(c *gin.Context) {
	err := ctrl.service.DeleteLoan(c.Request.Context(), c.Param("id"))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Loan deleted successfully!", nil)
}

func (ctrl *LoanController) PayFine(c *gin.Context) {
	id := c.Param("id")
	fine, err := ctrl.service.PayFine(c.Request.Context(), id)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Fine paid successfully!", resources.FormatFine(*fine))
}
