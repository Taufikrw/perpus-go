package controllers

import (
	"belajar-go/dto"
	"belajar-go/resources"
	"belajar-go/services"
	"belajar-go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookController struct {
	svc *services.BookService
}

func NewBookController(svc *services.BookService) *BookController {
	return &BookController{svc: svc}
}

func (ctrl *BookController) IndexBooks(c *gin.Context) {
	books, err := ctrl.svc.GetAllBooks(c.Request.Context())
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Books retrieved successfully", resources.FormatBooks(books))
}

func (ctrl *BookController) StoreBook(c *gin.Context) {
	var input dto.CreateBookDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		utils.HandleError(c, utils.NewValidationError("Invalid Input", errMsg))
		return
	}

	newBook, err := ctrl.svc.CreateBook(c.Request.Context(), input)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusCreated, "Book created successfully!", resources.FormatBook(*newBook))
}

func (ctrl *BookController) ShowBook(c *gin.Context) {
	book, err := ctrl.svc.GetBookByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Book detail retrieved successfully!", resources.FormatBook(*book))
}

func (ctrl *BookController) UpdateBook(c *gin.Context) {
	var input dto.CreateBookDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		utils.HandleError(c, utils.NewValidationError("Invalid Input", errMsg))
		return
	}

	updatedBook, err := ctrl.svc.UpdateBook(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Book updated successfully!", resources.FormatBook(*updatedBook))
}

func (ctrl *BookController) DeleteBook(c *gin.Context) {
	err := ctrl.svc.DeleteBook(c.Request.Context(), c.Param("id"))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Book deleted successfully!", nil)
}

func (ctrl *BookController) RestoreBook(c *gin.Context) {
	book, err := ctrl.svc.RestoreBook(c.Request.Context(), c.Param("id"))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Book restored successfully!", resources.FormatBook(*book))
}

func (ctrl *BookController) InsertBookItem(c *gin.Context) {
	var input dto.CreateBookItemDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		utils.HandleError(c, utils.NewValidationError("Invalid Input", errMsg))
		return
	}

	book, err := ctrl.svc.GetBookByID(c.Request.Context(), input.BookID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	newBookItem, err := ctrl.svc.CreateBookItem(c.Request.Context(), book.ID.String(), input)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusCreated, "Book item created successfully!", resources.FormatBookItem(*newBookItem))
}

func (ctrl *BookController) UpdateBookItem(c *gin.Context) {
	var input dto.UpdateBookItemDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		utils.HandleError(c, utils.NewValidationError("Invalid Input", errMsg))
		return
	}

	updatedBookItem, err := ctrl.svc.UpdateBookItem(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Book item updated successfully!", resources.FormatBookItem(*updatedBookItem))
}

func (ctrl *BookController) RemoveBookItem(c *gin.Context) {
	err := ctrl.svc.DeleteBookItem(c.Request.Context(), c.Param("id"))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Book item deleted successfully!", nil)
}

func (ctrl *BookController) RestoreBookItem(c *gin.Context) {
	item, err := ctrl.svc.RestoreBookItem(c.Request.Context(), c.Param("id"))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Book item restored successfully!", resources.FormatBookItem(*item))
}

func (ctrl *BookController) ShowBookItems(c *gin.Context) {
	bookItems, err := ctrl.svc.GetBookItemsByBookID(c.Request.Context(), c.Param("id"))
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Book items retrieved successfully!", resources.FormatBookItems(bookItems))
}
