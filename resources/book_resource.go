package resources

import "belajar-go/models"

type BookResource struct {
	ID        string           `json:"id"`
	Title     string           `json:"title"`
	Author    string           `json:"author"`
	Year      int              `json:"year"`
	Publisher string           `json:"publisher"`
	Isbn      string           `json:"isbn"`
	Synopsis  string           `json:"synopsis"`
	Category  CategoryResource `json:"category"`
}

type BookItemResource struct {
	ID            string       `json:"id"`
	InventoryCode string       `json:"inventory_code"`
	Condition     string       `json:"condition"`
	Status        string       `json:"status"`
	Book          BookResource `json:"book"`
}

func FormatBook(book models.Book) BookResource {
	return BookResource{
		ID:        book.ID.String(),
		Title:     book.Title,
		Author:    book.Author,
		Year:      book.Year,
		Publisher: book.Publisher,
		Isbn:      book.Isbn,
		Synopsis:  book.Synopsis,
		Category:  FormatCategory(book.Category),
	}
}

func FormatBooks(books []models.Book) []BookResource {
	var bookResources []BookResource
	for _, book := range books {
		bookResources = append(bookResources, FormatBook(book))
	}
	return bookResources
}

func FormatBookItem(bookItem models.BookItem) BookItemResource {
	return BookItemResource{
		ID:            bookItem.ID.String(),
		InventoryCode: bookItem.InventoryCode,
		Condition:     bookItem.Condition,
		Status:        bookItem.Status,
		Book:          FormatBook(bookItem.Book),
	}
}

func FormatBookItems(bookItems []models.BookItem) []BookItemResource {
	var bookItemResources []BookItemResource
	for _, bookItem := range bookItems {
		bookItemResources = append(bookItemResources, FormatBookItem(bookItem))
	}
	return bookItemResources
}
