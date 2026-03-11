package controllers

import (
	"belajar-go/dto"
	"belajar-go/resources"
	"belajar-go/services"
	"belajar-go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MemberController struct {
	svc *services.MemberService
}

func NewMemberController(svc *services.MemberService) *MemberController {
	return &MemberController{svc: svc}
}

func (ctrl *MemberController) IndexMember(c *gin.Context) {
	members, err := ctrl.svc.GetAllMembers(c.Request.Context())
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Members retrieved successfully!", resources.FormatMembers(members))
}

func (ctrl *MemberController) StoreMember(c *gin.Context) {
	var input dto.CreateMemberDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		utils.HandleError(c, utils.NewValidationError("Invalid Input", errMsg))
		return
	}

	newMember, err := ctrl.svc.CreateMember(c.Request.Context(), input)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusCreated, "Member created successfully!", resources.FormatMember(*newMember))
}

func (ctrl *MemberController) ShowMember(c *gin.Context) {
	member, err := ctrl.svc.GetMemberByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Member retrieved successfully!", resources.FormatMember(*member))
}

func (ctrl *MemberController) UpdateMember(c *gin.Context) {
	member, err := ctrl.svc.GetMemberByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	var input dto.UpdateMemberDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		utils.HandleError(c, utils.NewValidationError("Invalid Input", errMsg))
		return
	}

	updatedMember, err := ctrl.svc.UpdateMember(c.Request.Context(), member.ID.String(), input)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Member updated successfully!", resources.FormatMember(*updatedMember))
}

func (ctrl *MemberController) DeleteMember(c *gin.Context) {
	err := ctrl.svc.DeleteMember(c.Request.Context(), c.Param("id"))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Member deleted successfully!", nil)
}

func (ctrl *MemberController) RestoreMember(c *gin.Context) {
	member, err := ctrl.svc.RestoreMember(c.Request.Context(), c.Param("id"))
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Member restored successfully!", resources.FormatMember(*member))
}

func (ctrl *MemberController) ApproveMember(c *gin.Context) {
	member, err := ctrl.svc.ApproveMember(c.Request.Context(), c.Param("id"))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Member approved successfully!", resources.FormatMember(*member))
}
