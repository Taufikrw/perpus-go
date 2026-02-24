package routes

import (
	"belajar-go/controllers"
	"belajar-go/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")

	api.POST("/sign-in", controllers.Login)
	api.POST("/sign-up", controllers.Register)

	api.Use(middleware.RequireAuth())
	{
		api.GET("/books", controllers.IndexBooks)
		api.GET("/books/:id", controllers.ShowBook)
		api.GET("/books/:id/item", controllers.ShowBookItems)
		api.GET("/loans/:id", middleware.RequireLoanAccess(), controllers.ShowLoan)

		admin := api.Group("/")
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.GET("/categories", controllers.IndexCategory)
			admin.POST("/categories", controllers.StoreCategory)
			admin.GET("/categories/:id", controllers.ShowCategory)
			admin.PUT("/categories/:id", controllers.UpdateCategory)
			admin.DELETE("/categories/:id", controllers.DeleteCategory)

			admin.GET("/members", controllers.IndexMember)
			admin.POST("/members", controllers.StoreMember)
			admin.GET("/members/:id", controllers.ShowMember)
			admin.PUT("/members/:id", controllers.UpdateMember)
			admin.DELETE("/members/:id", controllers.DeleteMember)
		}

		adminOrLibrarian := api.Group("/")
		adminOrLibrarian.Use(middleware.RequireRole("admin", "librarian"))
		{
			adminOrLibrarian.POST("/books", controllers.StoreBook)
			adminOrLibrarian.PUT("/books/:id", controllers.UpdateBook)
			adminOrLibrarian.DELETE("/books/:id", controllers.DeleteBook)
			adminOrLibrarian.POST("/books/item", controllers.InsertBookItem)
			adminOrLibrarian.PUT("/books/item/:id", controllers.UpdateBookItem)
			adminOrLibrarian.DELETE("/books/item/:id", controllers.RemoveBookItem)

			adminOrLibrarian.GET("/loans", controllers.IndexLoans)
			adminOrLibrarian.PUT("/loans/:id", controllers.UpdateLoan)
			adminOrLibrarian.PUT("/loans/:id/return", controllers.ReturnLoan)
			adminOrLibrarian.DELETE("/loans/:id", controllers.DeleteLoan)
			adminOrLibrarian.GET("/fines", controllers.IndexFines)
			adminOrLibrarian.PUT("/loans/:id/pay-fine", controllers.PayFine)
			adminOrLibrarian.PUT("/approve/:id", controllers.ApproveMember)
		}

		member := api.Group("/")
		member.Use(middleware.RequireRole("member"))
		{
			member.POST("/loans", controllers.StoreLoan)
		}
	}
	return r
}
