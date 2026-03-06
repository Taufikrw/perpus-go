package routes

import (
	"belajar-go/controllers"
	"belajar-go/middleware"
	"belajar-go/repository"
	"belajar-go/services"

	"gorm.io/gorm"
)

type AppControllers struct {
	Category   *controllers.CategoryController
	Book       *controllers.BookController
	Auth       *controllers.AuthController
	Member     *controllers.MemberController
	Loan       *controllers.LoanController
	Fine       *controllers.FineController
	Middleware *middleware.AppMiddleware
}

func InitDependency(db *gorm.DB) *AppControllers {
	tx := repository.NewTransactionManager(db)

	categoryRepo := repository.NewCategoryRepository(db)
	bookRepo := repository.NewBookRepository(db)
	bookItemRepo := repository.NewBookItemRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	userRepo := repository.NewUserRepository(db)
	loanRepo := repository.NewLoanRepository(db)
	fineRepo := repository.NewFineRepository(db)

	categoryService := services.NewCategoryService(categoryRepo)
	bookService := services.NewBookService(bookRepo, categoryRepo, bookItemRepo)
	memberService := services.NewMemberService(tx, memberRepo, userRepo)
	loanService := services.NewLoanService(tx, loanRepo, memberRepo, bookItemRepo, fineRepo)
	authService := services.NewAuthService(userRepo, memberRepo, tx)
	fineService := services.NewFineService(fineRepo, loanRepo)

	appMiddlewre := middleware.NewAppMiddleware(db)

	return &AppControllers{
		Category:   controllers.NewCategoryController(categoryService),
		Book:       controllers.NewBookController(bookService),
		Auth:       controllers.NewAuthController(authService),
		Member:     controllers.NewMemberController(memberService),
		Loan:       controllers.NewLoanController(loanService),
		Fine:       controllers.NewFineController(fineService),
		Middleware: appMiddlewre,
	}
}
