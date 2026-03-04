package middleware

import (
	"belajar-go/models"
	"belajar-go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (m *AppMiddleware) RequireLoanAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.SendErrorResponse(c, 401, "Unauthorized", "Invalid token user not found")
			c.Abort()
			return
		}

		var user models.User
		if err := m.db.Preload("Role").Where("id = ?", userID).Take(&user).Error; err != nil {
			utils.SendErrorResponse(c, 404, "Not Found", "User not found")
			c.Abort()
			return
		}
		userRole := user.Role.Name
		id := c.Param("id")

		var loan models.Loan
		if err := m.db.Preload("Member.User.Role").Where("id = ?", id).Take(&loan).Error; err != nil {
			utils.SendErrorResponse(c, 404, "Not Found", "Loan not found")
			c.Abort()
			return
		}

		switch userRole {
		case "admin", "librarian":
			c.Next()
			return
		case "member":
			if loan.Member.UserID.String() != userID.(string) {
				utils.SendErrorResponse(c, http.StatusForbidden, "Access Denied", "You don't have permission to access this resource because you are not the owner of this loan")
				c.Abort()
				return
			}
			c.Next()
			return
		}

		utils.SendErrorResponse(c, http.StatusForbidden, "Access Denied", "You don't have permission to access this resource")
		c.Abort()
	}
}
