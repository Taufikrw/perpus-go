package resources

import (
	"belajar-go/models"
	"belajar-go/utils"
)

type FineResource struct {
	ID     string       `json:"id"`
	Loan   LoanResource `json:"loan_id"`
	Amount int          `json:"amount"`
	PaidAt *string      `json:"paid_at"`
}

func FormatFine(fine models.Fine) FineResource {
	var paidAt *string
	if fine.PaidAt != nil {
		formattedDate := utils.FormatDate(*fine.PaidAt)
		paidAt = &formattedDate
	}

	return FineResource{
		ID:     fine.ID.String(),
		Loan:   FormatLoan(fine.Loan),
		Amount: int(fine.Amount),
		PaidAt: paidAt,
	}
}

func FormatFines(fines []models.Fine) []FineResource {
	var fineResources []FineResource
	for _, fine := range fines {
		fineResources = append(fineResources, FormatFine(fine))
	}
	return fineResources
}
