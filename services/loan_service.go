package services

import (
	"belajar-go/dto"
	"belajar-go/models"
	"belajar-go/utils"
	"context"
	"time"
)

type LoanService struct {
	tx           models.TransactionManager
	loanRepo     models.LoanRepositoryInterface
	memberRepo   models.MemberRepositoryInterface
	bookItemRepo models.BookItemRepositoryInterface
	fineRepo     models.FineRepositoryInterface
}

func NewLoanService(tx models.TransactionManager, loanRepo models.LoanRepositoryInterface, memberRepo models.MemberRepositoryInterface, bookItemRepo models.BookItemRepositoryInterface, fineRepo models.FineRepositoryInterface) *LoanService {
	return &LoanService{
		tx:           tx,
		loanRepo:     loanRepo,
		memberRepo:   memberRepo,
		bookItemRepo: bookItemRepo,
		fineRepo:     fineRepo,
	}
}

func (s *LoanService) GetAllLoans(c context.Context) ([]models.Loan, error) {
	return s.loanRepo.FindAll(c)
}

func (s *LoanService) GetLoanByID(c context.Context, id string) (*models.Loan, error) {
	loan, err := s.loanRepo.FindByID(c, id)
	if loan == nil {
		return nil, utils.NewNotFoundError("Loan not found")
	} else if err != nil {
		return nil, err
	}
	return loan, nil
}

func (s *LoanService) CreateLoan(c context.Context, input dto.CreateLoanDTO, userID string) (*models.Loan, error) {
	member, err := s.memberRepo.FindByUserID(c, userID)
	if err != nil {
		return nil, utils.NewNotFoundError("Member not found")
	}
	bookItem, err := s.bookItemRepo.FindByID(c, input.BookItemID)
	if err != nil {
		return nil, utils.NewNotFoundError("Book item not found")
	}
	if bookItem.Status != "available" {
		return nil, utils.NewBadRequestError("Book item is not available for loan")
	}
	loanDate, _ := time.Parse("2006-01-02", input.LoanDate)
	dueDate, _ := time.Parse("2006-01-02", input.DueDate)
	newLoan := &models.Loan{
		MemberID:   member.ID,
		BookItemID: bookItem.ID,
		LoanDate:   loanDate,
		DueDate:    dueDate,
		Status:     "ongoing",
	}

	err = s.tx.WithTransaction(c, func(txCtx context.Context) error {
		if err := s.loanRepo.Create(txCtx, newLoan); err != nil {
			return err
		}
		bookItem.Status = "loaned"
		if err := s.bookItemRepo.Update(txCtx, bookItem); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.GetLoanByID(c, newLoan.ID.String())
}

func (s *LoanService) UpdateLoan(c context.Context, id string, input dto.UpdateLoanDTO) (*models.Loan, error) {
	loan, err := s.GetLoanByID(c, id)
	if err != nil {
		return nil, err
	}

	member, err := s.memberRepo.FindByID(c, input.MemberID)
	if err != nil {
		return nil, utils.NewNotFoundError("Member not found")
	}
	bookItem, err := s.bookItemRepo.FindByID(c, input.BookItemID)
	if err != nil {
		return nil, utils.NewNotFoundError("Book item not found")
	}

	loanDate, _ := time.Parse("2006-01-02", input.LoanDate)
	dueDate, _ := time.Parse("2006-01-02", input.DueDate)
	var returnDate *time.Time
	if input.ReturnDate != "" {
		parsedReturnDate, _ := time.Parse("2006-01-02", input.ReturnDate)
		returnDate = &parsedReturnDate
	}

	loan.MemberID = member.ID
	loan.BookItemID = bookItem.ID
	loan.LoanDate = loanDate
	loan.DueDate = dueDate
	loan.ReturnDate = returnDate
	loan.Status = input.Status

	err = s.loanRepo.Update(c, loan)
	if err != nil {
		return nil, err
	}
	return s.GetLoanByID(c, loan.ID.String())
}

func (s *LoanService) DeleteLoan(c context.Context, id string) error {
	loan, err := s.GetLoanByID(c, id)
	if err != nil {
		return err
	}
	return s.loanRepo.Delete(c, loan)
}

func (s *LoanService) ReturnLoan(c context.Context, id string) (*models.Loan, error) {
	loan, err := s.GetLoanByID(c, id)
	if err != nil {
		return nil, err
	}
	if loan.Status != "ongoing" {
		return nil, utils.NewValidationError("Loan cannot be returned", []string{"Only ongoing loans can be returned"})
	}

	now := time.Now()
	loan.ReturnDate = &now
	if now.After(loan.DueDate) {
		loan.Status = "overdue"
	} else {
		loan.Status = "returned"
	}

	err = s.tx.WithTransaction(c, func(txCtx context.Context) error {
		if err := s.loanRepo.Update(txCtx, loan); err != nil {
			return err
		}
		if loan.Status == "overdue" {
			newFine := &models.Fine{
				LoanID: loan.ID,
				Amount: 5000 * float64(int(now.Sub(loan.DueDate).Hours()/24)),
			}
			if err := s.fineRepo.Create(txCtx, newFine); err != nil {
				return err
			}
		}
		bookItem, err := s.bookItemRepo.FindByID(txCtx, loan.BookItemID.String())
		if err != nil {
			return err
		}
		bookItem.Status = "available"
		if err := s.bookItemRepo.Update(txCtx, bookItem); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.GetLoanByID(c, loan.ID.String())
}

func (s *LoanService) PayFine(c context.Context, id string) (*models.Fine, error) {
	fine, err := s.fineRepo.FindByID(c, id)
	if err != nil {
		return nil, utils.NewNotFoundError("Fine not found")
	}
	if fine.PaidAt != nil {
		return nil, utils.NewValidationError("Fine is already paid", nil)
	}

	now := time.Now()
	fine.PaidAt = &now

	err = s.fineRepo.Update(c, fine)
	if err != nil {
		return nil, err
	}
	return fine, nil
}

func (s *LoanService) GetAllFines(c context.Context) ([]models.Fine, error) {
	return s.fineRepo.FindAll(c)
}
