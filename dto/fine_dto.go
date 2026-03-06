package dto

type FineDTO struct {
	LoanID string  `json:"loan_id" binding:"required,uuid4"`
	Amount float64 `json:"amount" binding:"required,number,gt=0"`
	PaidAt string  `json:"paid_at" binding:"datetime=2006-01-02"`
}
