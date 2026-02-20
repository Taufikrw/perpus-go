package dto

type CreateMemberDTO struct {
	Username    string `json:"username" binding:"required,unique_username"`
	Email       string `json:"email" binding:"required,email,unique_email"`
	Password    string `json:"password" binding:"required,min=6"`
	MemberCode  string `json:"member_code" binding:"required,unique_member_code"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	IsApproved  bool   `json:"is_approved"`
}

type UpdateMemberDTO struct {
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"omitempty,min=6"`
	MemberCode  string `json:"member_code" binding:"required"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	IsApproved  bool   `json:"is_approved"`
}
