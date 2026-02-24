package resources

import (
	"belajar-go/models"
	"belajar-go/utils"
)

type LoanResource struct {
	ID         string           `json:"id"`
	Member     MemberResource   `json:"member_id"`
	BookItem   BookItemResource `json:"book_id"`
	LoanDate   string           `json:"loan_date"`
	DueDate    string           `json:"due_date"`
	ReturnDate *string          `json:"return_date"`
	Status     string           `json:"status"`
}

func FormatLoan(loan models.Loan) LoanResource {
	var returnDate *string
	if loan.ReturnDate != nil {
		formattedDate := utils.FormatDate(*loan.ReturnDate)
		returnDate = &formattedDate
	}

	return LoanResource{
		ID:         loan.ID.String(),
		Member:     FormatMember(loan.Member),
		BookItem:   FormatBookItem(loan.BookItem),
		LoanDate:   utils.FormatDate(loan.LoanDate),
		DueDate:    utils.FormatDate(loan.DueDate),
		ReturnDate: returnDate,
		Status:     loan.Status,
	}
}

func FormatLoans(loans []models.Loan) []LoanResource {
	var loanResources []LoanResource
	for _, loan := range loans {
		loanResources = append(loanResources, FormatLoan(loan))
	}
	return loanResources
}
