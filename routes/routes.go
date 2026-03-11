package routes

import (
	"belajar-go/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(ctrl *AppControllers) *gin.Engine {
	r := gin.Default()
	authMid := ctrl.Middleware
	api := r.Group("/api")

	api.POST("/sign-in", ctrl.Auth.Login)
	api.POST("/sign-up", ctrl.Auth.Register)

	api.Use(middleware.RequireAuth())
	{
		api.GET("/books", ctrl.Book.IndexBooks)
		api.GET("/books/:id", ctrl.Book.ShowBook)
		api.GET("/books/:id/item", ctrl.Book.ShowBookItems)
		api.GET("/loans/:id", authMid.RequireLoanAccess(), ctrl.Loan.ShowLoan)
		api.GET("/categories", ctrl.Category.IndexCategory)

		admin := api.Group("/")
		admin.Use(authMid.RequireRole("admin"))
		{
			admin.POST("/categories", ctrl.Category.StoreCategory)
			admin.GET("/categories/:id", ctrl.Category.ShowCategory)
			admin.PUT("/categories/:id", ctrl.Category.UpdateCategory)
			admin.DELETE("/categories/:id", ctrl.Category.DeleteCategory)
			admin.PUT("/categories/:id/restore", ctrl.Category.RestoreCategory)

			admin.GET("/members", ctrl.Member.IndexMember)
			admin.POST("/members", ctrl.Member.StoreMember)
			admin.GET("/members/:id", ctrl.Member.ShowMember)
			admin.PUT("/members/:id", ctrl.Member.UpdateMember)
			admin.DELETE("/members/:id", ctrl.Member.DeleteMember)
			admin.PUT("/members/:id/restore", ctrl.Member.RestoreMember)
		}

		adminOrLibrarian := api.Group("/")
		adminOrLibrarian.Use(authMid.RequireRole("admin", "librarian"))
		{
			adminOrLibrarian.POST("/books", ctrl.Book.StoreBook)
			adminOrLibrarian.PUT("/books/:id", ctrl.Book.UpdateBook)
			adminOrLibrarian.DELETE("/books/:id", ctrl.Book.DeleteBook)
			adminOrLibrarian.PUT("/books/:id/restore", ctrl.Book.RestoreBook)
			adminOrLibrarian.POST("/books/item", ctrl.Book.InsertBookItem)
			adminOrLibrarian.PUT("/books/item/:id", ctrl.Book.UpdateBookItem)
			adminOrLibrarian.DELETE("/books/item/:id", ctrl.Book.RemoveBookItem)
			adminOrLibrarian.PUT("/books/item/:id/restore", ctrl.Book.RestoreBookItem)

			adminOrLibrarian.GET("/loans", ctrl.Loan.IndexLoans)
			adminOrLibrarian.PUT("/loans/:id", ctrl.Loan.UpdateLoan)
			adminOrLibrarian.PUT("/loans/:id/return", ctrl.Loan.ReturnLoan)
			adminOrLibrarian.DELETE("/loans/:id", ctrl.Loan.DeleteLoan)
			adminOrLibrarian.PUT("/loans/:id/pay-fine", ctrl.Loan.PayFine)
			adminOrLibrarian.PUT("/approve/:id", ctrl.Member.ApproveMember)

			adminOrLibrarian.GET("/fines", ctrl.Fine.IndexFines)
			adminOrLibrarian.GET("/fines/:id", ctrl.Fine.ShowFine)
			adminOrLibrarian.POST("/fines", ctrl.Fine.StoreFine)
			adminOrLibrarian.PUT("/fines/:id", ctrl.Fine.UpdateFine)
			adminOrLibrarian.DELETE("/fines/:id", ctrl.Fine.DeleteFine)
		}

		member := api.Group("/")
		member.Use(authMid.RequireRole("member"))
		{
			member.POST("/loans", ctrl.Loan.StoreLoan)
		}
	}
	return r
}
