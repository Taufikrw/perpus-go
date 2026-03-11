package services

import (
	"belajar-go/dto"
	"belajar-go/models"
	"belajar-go/repository"
	"belajar-go/utils"
	"context"
	"time"
)

type FineService struct {
	fineRepo repository.FineRepository
	loanRepo repository.LoanRepository
}

func NewFineService(fineRepo repository.FineRepository, loanRepo repository.LoanRepository) *FineService {
	return &FineService{fineRepo: fineRepo, loanRepo: loanRepo}
}

func (s *FineService) GetAllFines(c context.Context) ([]models.Fine, error) {
	return s.fineRepo.GetAll(c, "Loan.Member.User.Role", "Loan.BookItem.Book.Category")
}

func (s *FineService) GetFineByID(c context.Context, id string) (*models.Fine, error) {
	fine, err := s.fineRepo.GetByID(c, id, "Loan.Member.User.Role", "Loan.BookItem.Book.Category")
	if fine == nil {
		return nil, utils.NewNotFoundError("Fine not found")
	} else if err != nil {
		return nil, err
	}
	return fine, nil
}

func (s *FineService) CreateFine(c context.Context, input dto.FineDTO) (*models.Fine, error) {
	loan, err := s.loanRepo.GetByID(c, input.LoanID, "Member.User.Role", "BookItem.Book.Category", "Fine")
	if err != nil {
		return nil, utils.NewNotFoundError("Loan not found")
	}

	paidAt, _ := time.Parse("2006-01-02", input.PaidAt)

	newFine := models.Fine{
		LoanID: loan.ID,
		Amount: input.Amount,
		PaidAt: &paidAt,
	}

	err = s.fineRepo.Create(c, &newFine)
	if err != nil {
		return nil, err
	}
	return s.GetFineByID(c, newFine.ID.String())
}

func (s *FineService) UpdateFine(c context.Context, id string, input dto.FineDTO) (*models.Fine, error) {
	fine, err := s.GetFineByID(c, id)
	if err != nil {
		return nil, err
	}

	loan, err := s.loanRepo.GetByID(c, input.LoanID, "Member.User.Role", "BookItem.Book.Category", "Fine")
	if err != nil {
		return nil, utils.NewNotFoundError("Loan not found")
	}

	paidAt, _ := time.Parse("2006-01-02", input.PaidAt)

	fine.LoanID = loan.ID
	fine.Amount = input.Amount
	fine.PaidAt = &paidAt

	err = s.fineRepo.Update(c, fine)
	if err != nil {
		return nil, err
	}
	return s.GetFineByID(c, fine.ID.String())
}

func (s *FineService) DeleteFine(c context.Context, id string) error {
	fine, err := s.GetFineByID(c, id)
	if err != nil {
		return err
	}
	return s.fineRepo.Delete(c, fine)
}
