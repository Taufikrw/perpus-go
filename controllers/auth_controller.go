package controllers

import (
	"belajar-go/config"
	"belajar-go/dto"
	"belajar-go/models"
	"belajar-go/resources"
	"belajar-go/utils"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginUser(input dto.LoginInput) (string, error) {
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return "", errors.New("email atau password salah")
	}

	if !utils.CheckPassword(user.Password, input.Password) {
		return "", errors.New("email atau password salah")
	}

	token, err := utils.GenerateToken(user.ID.String(), user.Email)
	log.Printf("Generated token for user %s: %s", user.Email, token)
	if err != nil {
		return "", errors.New("gagal membuat sesi login")
	}

	return token, nil
}

func Login(c *gin.Context) {
	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	token, err := LoginUser(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Login berhasil!",
		"token":   "Bearer " + token,
	})
}

func Register(c *gin.Context) {
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
		IsApproved:  false,
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
		"message": "Registrasi berhasil! Tunggu persetujuan admin untuk mengaktifkan akun Anda.",
		"data":    resources.FormatMember(newMember),
	})
}
