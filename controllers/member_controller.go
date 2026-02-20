package controllers

import (
	"belajar-go/config"
	"belajar-go/dto"
	"belajar-go/models"
	"belajar-go/resources"
	"belajar-go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IndexMember(c *gin.Context) {
	var members []models.Member
	config.DB.Preload("User.Role").Find(&members)
	utils.SendResponse(c, http.StatusOK, "Daftar member berhasil diambil!", resources.FormatMembers(members))
}

func StoreMember(c *gin.Context) {
	var input dto.CreateMemberDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		utils.SendErrorResponse(c, http.StatusUnprocessableEntity, "Input tidak valid", errMsg)
		return
	}
	var roleMember models.Role
	if err := config.DB.Where("name = ?", "member").Take(&roleMember).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal menemukan role member", nil)
		return
	}

	tx := config.DB.Begin()
	hashPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal meng-hash password", nil)
		return
	}

	newUser := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashPassword,
		RoleID:   roleMember.ID,
	}
	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal menyimpan user", nil)
		return
	}

	newMember := models.Member{
		MemberCode:  input.MemberCode,
		PhoneNumber: input.PhoneNumber,
		Address:     input.Address,
		IsApproved:  input.IsApproved,
		User:        newUser,
	}
	if err := tx.Create(&newMember).Error; err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal menyimpan member", nil)
		return
	}

	tx.Commit()
	config.DB.Preload("User.Role").Take(&newMember, newMember.ID)

	utils.SendResponse(c, http.StatusCreated, "Member berhasil dibuat!", resources.FormatMember(newMember))
}

func ShowMember(c *gin.Context) {
	var member models.Member

	if err := config.DB.Where("id = ?", c.Param("id")).Preload("User.Role").Take(&member).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Member tidak ditemukan!", nil)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Member berhasil diambil!", resources.FormatMember(member))
}

func UpdateMember(c *gin.Context) {
	var member models.Member
	if err := config.DB.Where("id = ?", c.Param("id")).Take(&member).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Member tidak ditemukan!", nil)
		return
	}

	var input dto.UpdateMemberDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		utils.SendErrorResponse(c, http.StatusUnprocessableEntity, "Input tidak valid", errMsg)
		return
	}

	if err := config.DB.Where("member_code = ? AND id != ?", input.MemberCode, member.ID).Take(&models.Member{}); err == nil {
		utils.SendErrorResponse(c, http.StatusUnprocessableEntity, "Input tidak valid", gin.H{"member_code": "Member code ini sudah digunakan oleh member lain"})
		return
	}
	var userCheck models.User
	if err := config.DB.Where("email = ? AND id != ?", input.Email, member.UserID).Take(&userCheck); err == nil {
		utils.SendErrorResponse(c, http.StatusUnprocessableEntity, "Input tidak valid", gin.H{"email": "Email ini sudah digunakan oleh user lain"})
		return
	}
	if err := config.DB.Where("username = ? AND id != ?", input.Username, member.UserID).Take(&userCheck); err == nil {
		utils.SendErrorResponse(c, http.StatusUnprocessableEntity, "Input tidak valid", gin.H{"username": "Username ini sudah digunakan oleh user lain"})
		return
	}

	tx := config.DB.Begin()

	hashPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal meng-hash password", nil)
		return
	}

	var roleMember models.Role
	if err := config.DB.Where("name = ?", "member").Take(&roleMember).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal menemukan role member", nil)
		return
	}

	newUser := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashPassword,
		RoleID:   roleMember.ID,
	}
	if err := tx.Model(&models.User{}).Where("id = ?", member.UserID).Updates(&newUser).Error; err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengupdate user", nil)
		return
	}

	newMember := models.Member{
		MemberCode:  input.MemberCode,
		PhoneNumber: input.PhoneNumber,
		Address:     input.Address,
		IsApproved:  input.IsApproved,
	}
	if err := tx.Model(&member).Updates(&newMember).Error; err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengupdate member", nil)
		return
	}

	tx.Commit()
	config.DB.Preload("User.Role").Take(&member, member.ID)

	utils.SendResponse(c, http.StatusOK, "Member berhasil diupdate!", resources.FormatMember(member))
}

func DeleteMember(c *gin.Context) {
	var member models.Member
	if err := config.DB.Where("id = ?", c.Param("id")).Take(&member).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Member tidak ditemukan!", nil)
		return
	}

	tx := config.DB.Begin()

	if err := tx.Where("id = ?", member.ID).Delete(&member).Error; err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus member", nil)
		return
	}

	if err := tx.Where("id = ?", member.UserID).Delete(&member.User).Error; err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus user", nil)
		return
	}

	tx.Commit()
	utils.SendResponse(c, http.StatusOK, "Member berhasil dihapus!", nil)
}

func ApproveMember(c *gin.Context) {
	var member models.Member
	if err := config.DB.Where("id = ?", c.Param("id")).Take(&member).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Member tidak ditemukan!", nil)
		return
	}

	member.IsApproved = true
	if err := config.DB.Save(&member).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengupdate status approval member", nil)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Member berhasil disetujui!", resources.FormatMember(member))
}
