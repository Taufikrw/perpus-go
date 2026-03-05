package services

import (
	"belajar-go/dto"
	"belajar-go/models"
	"belajar-go/utils"
	"context"
)

type BookService struct {
	bookRepo     models.BookRepository
	categoryRepo models.CategoryRepository
	bookItemRepo models.BookItemRepository
}

func NewBookService(bookRepo models.BookRepository, categoryRepo models.CategoryRepository, bookItemRepo models.BookItemRepository) *BookService {
	return &BookService{bookRepo: bookRepo, categoryRepo: categoryRepo, bookItemRepo: bookItemRepo}
}

func (s *BookService) GetAllBooks(c context.Context) ([]models.Book, error) {
	return s.bookRepo.FindAll(c)
}

func (s *BookService) GetBookByID(c context.Context, id string) (*models.Book, error) {
	book, err := s.bookRepo.FindByID(c, id)
	if book == nil {
		return nil, utils.NewNotFoundError("Book not found")
	} else if err != nil {
		return nil, err
	}
	return book, nil
}

func (s *BookService) CreateBook(c context.Context, input dto.CreateBookDTO) (*models.Book, error) {
	category, err := s.categoryRepo.FindByID(c, input.CategoryID)
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

	category, err := s.categoryRepo.FindByID(c, input.CategoryID)
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

func (s *BookService) GetBookItemsByBookID(c context.Context, bookID string) ([]models.BookItem, error) {
	bookItems, err := s.bookItemRepo.FindByBookID(c, bookID)
	if err != nil {
		return nil, utils.NewNotFoundError("Book items not found for the given book ID")
	}
	return bookItems, nil
}

func (s *BookService) GetBookItemByID(c context.Context, id string) (*models.BookItem, error) {
	bookItem, err := s.bookItemRepo.FindByID(c, id)
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
		BookID:        book.ID,
		InventoryCode: input.InventoryCode,
		Condition:     input.Condition,
		Status:        "available",
	}

	err = s.bookItemRepo.Create(c, &newBookItem)
	if err != nil {
		return nil, err
	}
	return s.bookItemRepo.FindByID(c, newBookItem.ID.String())
}

func (s *BookService) UpdateBookItem(c context.Context, id string, input dto.UpdateBookItemDTO) (*models.BookItem, error) {
	bookItem, err := s.GetBookItemByID(c, id)
	if err != nil {
		return nil, err
	}

	book, err := s.GetBookByID(c, input.BookID)
	if err != nil {
		return nil, err
	}

	bookItem.BookID = book.ID
	bookItem.InventoryCode = input.InventoryCode
	bookItem.Condition = input.Condition

	err = s.bookItemRepo.Update(c, bookItem)
	if err != nil {
		return nil, err
	}
	return s.bookItemRepo.FindByID(c, bookItem.ID.String())
}

func (s *BookService) DeleteBookItem(c context.Context, id string) error {
	bookItem, err := s.GetBookItemByID(c, id)
	if err != nil {
		return err
	}
	return s.bookItemRepo.Delete(c, bookItem)
}
