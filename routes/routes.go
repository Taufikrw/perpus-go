package routes

import (
	"belajar-go/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/books", controllers.IndexBooks)
	r.POST("/books", controllers.StoreBook)
	r.GET("/books/:id", controllers.ShowBook)
	r.PUT("/books/:id", controllers.UpdateBook)
	r.DELETE("/books/:id", controllers.DeleteBook)

	r.GET("/books/:id/item", controllers.ShowBookItems)
	r.POST("/books/item", controllers.InsertBookItem)
	r.PUT("/books/item/:id", controllers.UpdateBookItem)
	r.DELETE("/books/item/:id", controllers.RemoveBookItem)

	r.GET("/categories", controllers.IndexCategory)
	r.POST("/categories", controllers.StoreCategory)
	r.GET("/categories/:id", controllers.ShowCategory)
	r.PUT("/categories/:id", controllers.UpdateCategory)
	r.DELETE("/categories/:id", controllers.DeleteCategory)

	r.GET("/members", controllers.IndexMember)
	r.POST("/members", controllers.StoreMember)
	r.GET("/members/:id", controllers.ShowMember)
	r.PUT("/members/:id", controllers.UpdateMember)
	r.DELETE("/members/:id", controllers.DeleteMember)

	r.POST("/sign-in", controllers.Login)
	r.POST("/sign-up", controllers.Register)
	r.PUT("/approve/:id", controllers.ApproveMember)
	return r
}
