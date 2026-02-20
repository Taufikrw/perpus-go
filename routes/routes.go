package routes

import (
	"belajar-go/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/book", controllers.IndexBuku)
	r.POST("/book", controllers.StoreBuku)
	r.GET("/book/:id", controllers.ShowBuku)
	r.PUT("/book/:id", controllers.UpdateBuku)
	r.DELETE("/book/:id", controllers.DeleteBuku)

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
