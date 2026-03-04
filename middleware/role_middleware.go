package middleware

import (
	"belajar-go/models"
	"belajar-go/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppMiddleware struct {
	db *gorm.DB
}

// NewAppMiddleware adalah constructor untuk inisialisasi
func NewAppMiddleware(db *gorm.DB) *AppMiddleware {
	return &AppMiddleware{db: db}
}

func (m *AppMiddleware) RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.SendErrorResponse(c, 401, "Unauthorized", "Invalid token user not found")
			return
		}

		var user models.User
		if err := m.db.Preload("Role").Where("id = ?", userID).Take(&user).Error; err != nil {
			utils.SendErrorResponse(c, 401, "Unauthorized", "User not found")
			return
		}

		userRole := user.Role.Name
		isAllowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			utils.SendErrorResponse(c, 403, "Forbidden", "You don't have permission to access this resource")
			c.Abort()
			return
		}

		c.Next()
	}
}
