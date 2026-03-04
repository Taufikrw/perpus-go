package routes

import (
	"belajar-go/controllers"
	"belajar-go/middleware"
	"belajar-go/repository"
	"belajar-go/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	txManager := repository.NewTransactionManager(db)
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryCtrl := controllers.NewCategoryController(categoryService)

	authRepo := repository.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo)
	authCtrl := controllers.NewAuthController(authService)

	bookItemRepo := repository.NewBookItemRepository(db)
	bookRepo := repository.NewBookRepository(db)
	bookService := services.NewBookService(bookRepo, categoryRepo, bookItemRepo)
	bookCtrl := controllers.NewBookController(bookService)

	userRepo := repository.NewUserRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	memberService := services.NewMemberService(txManager, memberRepo, userRepo)
	memberCtrl := controllers.NewMemberController(memberService)

	fineRepo := repository.NewFineRepository(db)
	loanRepo := repository.NewLoanRepository(db)
	loanService := services.NewLoanService(txManager, loanRepo, memberRepo, bookItemRepo, fineRepo)
	loanCtrl := controllers.NewLoanController(loanService)

	r := gin.Default()
	authMid := middleware.NewAppMiddleware(db)
	api := r.Group("/api")

	api.POST("/sign-in", authCtrl.Login)
	api.POST("/sign-up", authCtrl.Register)

	api.Use(middleware.RequireAuth())
	{
		api.GET("/books", bookCtrl.IndexBooks)
		api.GET("/books/:id", bookCtrl.ShowBook)
		api.GET("/books/:id/item", bookCtrl.ShowBookItems)
		api.GET("/loans/:id", authMid.RequireLoanAccess(), loanCtrl.ShowLoan)
		api.GET("/categories", categoryCtrl.IndexCategory)

		admin := api.Group("/")
		admin.Use(authMid.RequireRole("admin"))
		{
			admin.POST("/categories", categoryCtrl.StoreCategory)
			admin.GET("/categories/:id", categoryCtrl.ShowCategory)
			admin.PUT("/categories/:id", categoryCtrl.UpdateCategory)
			admin.DELETE("/categories/:id", categoryCtrl.DeleteCategory)

			admin.GET("/members", memberCtrl.IndexMember)
			admin.POST("/members", memberCtrl.StoreMember)
			admin.GET("/members/:id", memberCtrl.ShowMember)
			admin.PUT("/members/:id", memberCtrl.UpdateMember)
			admin.DELETE("/members/:id", memberCtrl.DeleteMember)
		}

		adminOrLibrarian := api.Group("/")
		adminOrLibrarian.Use(authMid.RequireRole("admin", "librarian"))
		{
			adminOrLibrarian.POST("/books", bookCtrl.StoreBook)
			adminOrLibrarian.PUT("/books/:id", bookCtrl.UpdateBook)
			adminOrLibrarian.DELETE("/books/:id", bookCtrl.DeleteBook)
			adminOrLibrarian.POST("/books/item", bookCtrl.InsertBookItem)
			adminOrLibrarian.PUT("/books/item/:id", bookCtrl.UpdateBookItem)
			adminOrLibrarian.DELETE("/books/item/:id", bookCtrl.RemoveBookItem)

			adminOrLibrarian.GET("/loans", loanCtrl.IndexLoans)
			adminOrLibrarian.PUT("/loans/:id", loanCtrl.UpdateLoan)
			adminOrLibrarian.PUT("/loans/:id/return", loanCtrl.ReturnLoan)
			adminOrLibrarian.DELETE("/loans/:id", loanCtrl.DeleteLoan)
			adminOrLibrarian.GET("/fines", loanCtrl.IndexFines)
			adminOrLibrarian.PUT("/loans/:id/pay-fine", loanCtrl.PayFine)
			adminOrLibrarian.PUT("/approve/:id", memberCtrl.ApproveMember)
		}

		member := api.Group("/")
		member.Use(authMid.RequireRole("member"))
		{
			member.POST("/loans", loanCtrl.StoreLoan)
		}
	}
	return r
}
