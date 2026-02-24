package dto

type CreateLoanDTO struct {
	MemberID   string `json:"member_id" binding:"required,uuid4"`
	BookItemID string `json:"book_item_id" binding:"required,uuid4"`
	LoanDate   string `json:"loan_date" binding:"required,datetime=2006-01-02"`
	DueDate    string `json:"due_date" binding:"required,datetime=2006-01-02"`
}

type UpdateLoanDTO struct {
	MemberID   string `json:"member_id" binding:"required,uuid4"`
	BookItemID string `json:"book_item_id" binding:"required,uuid4"`
	LoanDate   string `json:"loan_date" binding:"required,datetime=2006-01-02"`
	DueDate    string `json:"due_date" binding:"required,datetime=2006-01-02"`
	ReturnDate string `json:"return_date" binding:"omitempty,datetime=2006-01-02"`
	Status     string `json:"status" binding:"required,oneof=ongoing returned overdue"`
}
