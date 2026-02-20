package controllers

import (
	"belajar-go/config"
	"belajar-go/dto"
	"belajar-go/models"
	"belajar-go/resources"
	"belajar-go/utils"

	"github.com/gin-gonic/gin"
)

func IndexCategory(c *gin.Context) {
	var categories []models.BookCategory
	config.DB.Find(&categories)
	utils.SendResponse(c, 200, "Daftar kategori berhasil diambil!", resources.FormatCategories(categories))
}

func StoreCategory(c *gin.Context) {
	var input dto.CategoryDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		utils.SendErrorResponse(c, 422, "Input tidak valid", errMsg)
		return
	}

	newCategory := models.BookCategory{
		Name: input.Name,
	}

	config.DB.Create(&newCategory)
	utils.SendResponse(c, 201, "Kategori berhasil dibuat!", resources.FormatCategory(newCategory))
}

func ShowCategory(c *gin.Context) {
	var category models.BookCategory
	if err := config.DB.Where("id = ?", c.Param("id")).Take(&category).Error; err != nil {
		utils.SendErrorResponse(c, 404, "Kategori tidak ditemukan!", nil)
		return
	}
	utils.SendResponse(c, 200, "Detail kategori berhasil diambil!", resources.FormatCategory(category))
}

func UpdateCategory(c *gin.Context) {
	var category models.BookCategory
	if err := config.DB.Where("id = ?", c.Param("id")).Take(&category).Error; err != nil {
		utils.SendErrorResponse(c, 404, "Kategori tidak ditemukan!", nil)
		return
	}

	var input dto.CategoryDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		utils.SendErrorResponse(c, 422, "Input tidak valid", errMsg)
		return
	}

	category.Name = input.Name
	config.DB.Updates(&category)
	utils.SendResponse(c, 200, "Kategori berhasil diupdate!", resources.FormatCategory(category))
}

func DeleteCategory(c *gin.Context) {
	var category models.BookCategory
	if err := config.DB.Where("id = ?", c.Param("id")).Take(&category).Error; err != nil {
		utils.SendErrorResponse(c, 404, "Kategori tidak ditemukan!", nil)
		return
	}

	config.DB.Delete(&category)
	utils.SendResponse(c, 200, "Kategori berhasil dihapus!", nil)
}
