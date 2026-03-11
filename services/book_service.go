package services

import (
	"belajar-go/dto"
	"belajar-go/models"
	"belajar-go/repository"
	"belajar-go/utils"
	"context"
)

type BookService struct {
	bookRepo     repository.BookRepository
	categoryRepo repository.CategoryRepository
	bookItemRepo repository.BookItemRepository
}

func NewBookService(bookRepo repository.BookRepository, categoryRepo repository.CategoryRepository, bookItemRepo repository.BookItemRepository) *BookService {
	return &BookService{bookRepo: bookRepo, categoryRepo: categoryRepo, bookItemRepo: bookItemRepo}
}

func (s *BookService) GetAllBooks(c context.Context) ([]models.Book, error) {
	return s.bookRepo.GetAll(c, "Category")
}

func (s *BookService) GetBookByID(c context.Context, id string) (*models.Book, error) {
	book, err := s.bookRepo.GetByID(c, id, "Category")
	if book == nil {
		return nil, utils.NewNotFoundError("Book not found")
	} else if err != nil {
		return nil, err
	}
	return book, nil
}

func (s *BookService) CreateBook(c context.Context, input dto.CreateBookDTO) (*models.Book, error) {
	category, err := s.categoryRepo.GetByID(c, input.CategoryID)
	if err != nil {
		return nil, utils.NewNotFoundError("Category not found")
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

	err = s.bookRepo.Create(c, &newBook)
	if err != nil {
		return nil, err
	}
	return s.GetBookByID(c, newBook.ID.String())
}

func (s *BookService) UpdateBook(c context.Context, id string, input dto.CreateBookDTO) (*models.Book, error) {
	book, err := s.GetBookByID(c, id)
	if err != nil {
		return nil, err
	}

	category, err := s.categoryRepo.GetByID(c, input.CategoryID)
	if err != nil {
		return nil, utils.NewNotFoundError("Category not found")
	}

	book.Title = input.Title
	book.Author = input.Author
	book.Year = input.Year
	book.Publisher = input.Publisher
	book.Isbn = input.Isbn
	book.Synopsis = input.Synopsis
	book.CategoryID = category.ID

	err = s.bookRepo.Update(c, book)
	if err != nil {
		return nil, err
	}
	return s.GetBookByID(c, book.ID.String())
}

func (s *BookService) DeleteBook(c context.Context, id string) error {
	book, err := s.GetBookByID(c, id)
	if err != nil {
		return err
	}
	return s.bookRepo.Delete(c, book)
}

func (s *BookService) RestoreBook(c context.Context, id string) (*models.Book, error) {
	book, _ := s.GetBookByID(c, id)
	if book != nil {
		return nil, utils.NewBadRequestError("Book is not deleted")
	}
	err := s.bookRepo.Restore(c, id)
	if err != nil {
		return nil, err
	}

	return s.GetBookByID(c, id)
}

func (s *BookService) GetBookItemsByBookID(c context.Context, bookID string) ([]models.BookItem, error) {
	bookItems, err := s.bookItemRepo.FindByBookID(c, bookID)
	if err != nil {
		return nil, utils.NewNotFoundError("Book items not found for the given book ID")
	}
	return bookItems, nil
}

func (s *BookService) GetBookItemByID(c context.Context, id string) (*models.BookItem, error) {
	bookItem, err := s.bookItemRepo.GetByID(c, id, "Book.Category")
	if bookItem == nil {
		return nil, utils.NewNotFoundError("Book item not found")
	} else if err != nil {
		return nil, err
	}
	return bookItem, nil
}

func (s *BookService) CreateBookItem(c context.Context, bookID string, input dto.CreateBookItemDTO) (*models.BookItem, error) {
	book, err := s.GetBookByID(c, bookID)
	if err != nil {
		return nil, err
	}

	newBookItem := models.BookItem{
		BookID:        &book.ID,
		InventoryCode: input.InventoryCode,
		Condition:     input.Condition,
		Status:        "available",
	}

	err = s.bookItemRepo.Create(c, &newBookItem)
	if err != nil {
		return nil, err
	}
	return s.bookItemRepo.GetByID(c, newBookItem.ID.String(), "Book.Category")
}

func (s *BookService) UpdateBookItem(c context.Context, id string, input dto.UpdateBookItemDTO) (*models.BookItem, error) {
	bookItem, err := s.GetBookItemByID(c, id)
	if err != nil {
		return nil, err
	}

	exist, err := s.bookItemRepo.IsInventoryCodeExists(c, input.InventoryCode, id)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, utils.NewValidationError("Validation failed", []string{"Inventory code already exists"})
	}

	book, err := s.GetBookByID(c, input.BookID)
	if err != nil {
		return nil, err
	}

	bookItem.BookID = &book.ID
	bookItem.InventoryCode = input.InventoryCode
	bookItem.Condition = input.Condition

	err = s.bookItemRepo.Update(c, bookItem)
	if err != nil {
		return nil, err
	}
	return s.bookItemRepo.GetByID(c, bookItem.ID.String(), "Book.Category")
}

func (s *BookService) DeleteBookItem(c context.Context, id string) error {
	bookItem, err := s.GetBookItemByID(c, id)
	if err != nil {
		return err
	}
	return s.bookItemRepo.Delete(c, bookItem)
}

func (s *BookService) RestoreBookItem(c context.Context, id string) (*models.BookItem, error) {
	item, _ := s.bookItemRepo.GetByID(c, id)
	if item != nil {
		return nil, utils.NewBadRequestError("Book item is not deleted")
	}
	err := s.bookItemRepo.Restore(c, id)
	if err != nil {
		return nil, err
	}

	return s.bookItemRepo.GetByID(c, id, "Book.Category")
}
