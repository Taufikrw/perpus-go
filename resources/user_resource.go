package resources

import "belajar-go/models"

type UserResource struct {
	ID       string       `json:"id"`
	Username string       `json:"username"`
	Email    string       `json:"email"`
	Role     RoleResource `json:"role"`
}

type RoleResource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func FormatUser(user models.User) UserResource {
	return UserResource{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Role: RoleResource{
			ID:   user.Role.ID.String(),
			Name: user.Role.Name,
		},
	}
}
