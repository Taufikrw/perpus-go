package controllers

import (
	"belajar-go/dto"
	"belajar-go/resources"
	"belajar-go/services"
	"belajar-go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	svc *services.AuthService
}

func NewAuthController(svc *services.AuthService) *AuthController {
	return &AuthController{svc: svc}
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		utils.HandleError(c, utils.NewValidationError("Invalid Input", errMsg))
		return
	}

	token, err := ctrl.svc.Login(c, input)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Login successful!", gin.H{"token": "Bearer " + token})
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var input dto.CreateMemberDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		utils.HandleError(c, utils.NewValidationError("Invalid Input", errMsg))
		return
	}

	newMember, err := ctrl.svc.Register(c.Request.Context(), input)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.SendResponse(c, http.StatusCreated, "Registration successful! Please wait for admin approval.", resources.FormatMember(*newMember))
}
