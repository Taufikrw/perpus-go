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

func IndexBooks(c *gin.Context) {
	var books []models.Book
	config.DB.Preload("Category").Find(&books)
	utils.SendResponse(c, http.StatusOK, "Daftar buku berhasil diambil!", resources.FormatBooks(books))
}

func StoreBook(c *gin.Context) {
	var input dto.CreateBookDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		utils.SendErrorResponse(c, http.StatusUnprocessableEntity, "Input tidak valid", errMsg)
		return
	}
	var category models.BookCategory
	if err := config.DB.Where("id = ?", input.CategoryID).Take(&category).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Kategori tidak ditemukan", nil)
		return
	}

	newBook := models.Book{
		Title:      input.Title,
		Author:     input.Author,
		Year:       input.Year,
		Publisher:  input.Publisher,
		Isbn:       input.Isbn,
		Synopsis:   input.Synopsis,
		CategoryID: category.ID,
	}

	config.DB.Create(&newBook)
	config.DB.Preload("Category").Take(&newBook, newBook.ID)
	utils.SendResponse(c, http.StatusCreated, "Buku berhasil dibuat!", resources.FormatBook(newBook))
}

func ShowBook(c *gin.Context) {
	var buku models.Book

	if err := config.DB.Where("id = ?", c.Param("id")).Preload("Category").Take(&buku).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Buku tidak ditemukan!", nil)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Detail buku berhasil diambil!", resources.FormatBook(buku))
}

func UpdateBook(c *gin.Context) {
	var buku models.Book
	if err := config.DB.Where("id = ?", c.Param("id")).Take(&buku).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Buku tidak ditemukan!", nil)
		return
	}

	var input dto.CreateBookDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		utils.SendErrorResponse(c, http.StatusUnprocessableEntity, "Input tidak valid", errMsg)
		return
	}
	var category models.BookCategory
	if err := config.DB.Where("id = ?", input.CategoryID).Take(&category).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Kategori tidak ditemukan", nil)
		return
	}

	newBook := models.Book{
		Title:      input.Title,
		Author:     input.Author,
		Year:       input.Year,
		Publisher:  input.Publisher,
		Isbn:       input.Isbn,
		Synopsis:   input.Synopsis,
		CategoryID: category.ID,
	}

	config.DB.Model(&buku).Updates(&newBook)
	config.DB.Preload("Category").Take(&newBook, buku.ID)
	utils.SendResponse(c, http.StatusOK, "Buku berhasil diupdate!", resources.FormatBook(newBook))
}

func DeleteBook(c *gin.Context) {
	var buku models.Book
	if err := config.DB.Where("id = ?", c.Param("id")).Take(&buku).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Buku tidak ditemukan!", nil)
		return
	}

	config.DB.Delete(&buku)
	utils.SendResponse(c, http.StatusOK, "Buku berhasil dihapus!", nil)
}
