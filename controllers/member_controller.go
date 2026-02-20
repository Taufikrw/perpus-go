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
	c.JSON(http.StatusOK, gin.H{"data": resources.FormatMembers(members)})
}

func StoreMember(c *gin.Context) {
	var input dto.CreateMemberDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}
	var roleMember models.Role
	if err := config.DB.Where("name = ?", "member").Take(&roleMember).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menemukan role member"})
		return
	}

	tx := config.DB.Begin()
	hashPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal meng-hash password"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan user"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan member"})
		return
	}

	tx.Commit()
	config.DB.Preload("User.Role").Take(&newMember, newMember.ID)

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Member berhasil dibuat!",
		"data":    resources.FormatMember(newMember),
	})
}

func ShowMember(c *gin.Context) {
	var member models.Member

	if err := config.DB.Where("id = ?", c.Param("id")).Preload("User.Role").Take(&member).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member tidak ditemukan!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resources.FormatMember(member)})
}

func UpdateMember(c *gin.Context) {
	var member models.Member
	if err := config.DB.Where("id = ?", c.Param("id")).Take(&member).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member tidak ditemukan!"})
		return
	}

	var input dto.UpdateMemberDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	if err := config.DB.Where("member_code = ? AND id != ?", input.MemberCode, member.ID).Take(&models.Member{}); err == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": gin.H{"member_code": "Kode member ini sudah digunakan oleh member lain"},
		})
		return
	}
	var userCheck models.User
	if err := config.DB.Where("email = ? AND id != ?", input.Email, member.UserID).Take(&userCheck); err == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": gin.H{"email": "Email ini sudah digunakan oleh user lain"},
		})
		return
	}
	if err := config.DB.Where("username = ? AND id != ?", input.Username, member.UserID).Take(&userCheck); err == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": gin.H{"username": "Username ini sudah digunakan oleh user lain"},
		})
		return
	}

	tx := config.DB.Begin()

	hashPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal meng-hash password"})
		return
	}

	var roleMember models.Role
	if err := config.DB.Where("name = ?", "member").Take(&roleMember).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menemukan role member"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate user"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate member"})
		return
	}

	tx.Commit()
	config.DB.Preload("User.Role").Take(&member, member.ID)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Member berhasil diupdate!",
		"data":    resources.FormatMember(member),
	})
}

func DeleteMember(c *gin.Context) {
	var member models.Member
	if err := config.DB.Where("id = ?", c.Param("id")).Take(&member).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member tidak ditemukan!"})
		return
	}

	tx := config.DB.Begin()

	if err := tx.Delete(&member.User).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus user"})
		return
	}

	if err := tx.Delete(&member).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus member"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Member berhasil dihapus!",
	})
}

func ApproveMember(c *gin.Context) {
	var member models.Member
	if err := config.DB.Where("id = ?", c.Param("id")).Take(&member).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member tidak ditemukan!"})
		return
	}

	member.IsApproved = true
	if err := config.DB.Save(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate status approval member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Status approval member berhasil diupdate!",
		"data":    resources.FormatMember(member),
	})
}
