package services

import (
	"belajar-go/dto"
	"belajar-go/models"
	"belajar-go/utils"
	"context"
)

type MemberService struct {
	tx         models.TransactionManager
	memberRepo models.MemberRepositoryInterface
	userRepo   models.UserRepositoryInterface
}

func NewMemberService(tx models.TransactionManager, memberRepo models.MemberRepositoryInterface, userRepo models.UserRepositoryInterface) *MemberService {
	return &MemberService{tx: tx, memberRepo: memberRepo, userRepo: userRepo}
}

func (s *MemberService) GetAllMembers(c context.Context) ([]models.Member, error) {
	return s.memberRepo.FindAll(c)
}

func (s *MemberService) GetMemberByID(c context.Context, id string) (*models.Member, error) {
	member, err := s.memberRepo.FindByID(c, id)
	if member == nil {
		return nil, utils.NewNotFoundError("Member not found")
	} else if err != nil {
		return nil, err
	}
	return member, nil
}

func (s *MemberService) CreateMember(c context.Context, input dto.CreateMemberDTO) (*models.Member, error) {
	roleMember, err := s.userRepo.GetRoleByName(c, "member")
	if err != nil {
		return nil, utils.NewNotFoundError("Role 'member' not found")
	}

	hashPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	newUser := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashPassword,
		RoleID:   roleMember.ID,
	}

	newMember := models.Member{
		MemberCode:  input.MemberCode,
		PhoneNumber: input.PhoneNumber,
		Address:     input.Address,
		IsApproved:  input.IsApproved,
		UserID:      newUser.ID,
	}

	err = s.tx.WithTransaction(c, func(txCtx context.Context) error {
		if err := s.userRepo.Create(txCtx, &newUser); err != nil {
			return err
		}
		newMember.UserID = newUser.ID

		if err := s.memberRepo.Create(txCtx, &newMember); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.GetMemberByID(c, newMember.ID.String())
}

func (s *MemberService) UpdateMember(c context.Context, id string, input dto.UpdateMemberDTO) (*models.Member, error) {
	member, err := s.GetMemberByID(c, id)
	if err != nil {
		return nil, err
	}

	validationErr := make(map[string]string)
	exists, err := s.memberRepo.IsMemberCodeExists(c, input.MemberCode, member.ID.String())
	if err != nil {
		return nil, err
	}
	if exists {
		validationErr["member_code"] = "Member code already exists"
	}
	exists, err = s.userRepo.IsEmailExists(c, input.Email, member.UserID.String())
	if err != nil {
		return nil, err
	}
	if exists {
		validationErr["email"] = "Email already exists"
	}
	exists, err = s.userRepo.IsUsernameExists(c, input.Username, member.UserID.String())
	if err != nil {
		return nil, err
	}
	if exists {
		validationErr["username"] = "Username already exists"
	}
	if len(validationErr) > 0 {
		return nil, utils.NewValidationError("Validation failed", validationErr)
	}

	hashPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	member.MemberCode = input.MemberCode
	member.PhoneNumber = input.PhoneNumber
	member.Address = input.Address
	member.IsApproved = input.IsApproved

	member.User.Username = input.Username
	member.User.Email = input.Email
	member.User.Password = hashPassword

	err = s.tx.WithTransaction(c, func(txCtx context.Context) error {
		if err := s.userRepo.Update(txCtx, &member.User); err != nil {
			return err
		}
		if err := s.memberRepo.Update(txCtx, member); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.GetMemberByID(c, member.ID.String())
}

func (s *MemberService) DeleteMember(c context.Context, id string) error {
	member, err := s.GetMemberByID(c, id)
	if err != nil {
		return err
	}

	return s.tx.WithTransaction(c, func(txCtx context.Context) error {
		if err := s.memberRepo.Delete(txCtx, member); err != nil {
			return err
		}
		if err := s.userRepo.Delete(txCtx, &member.User); err != nil {
			return err
		}
		return nil
	})
}

func (s *MemberService) ApproveMember(c context.Context, id string) (*models.Member, error) {
	member, err := s.GetMemberByID(c, id)
	if err != nil {
		return nil, err
	}

	member.IsApproved = true
	err = s.memberRepo.Update(c, member)
	if err != nil {
		return nil, err
	}
	return s.GetMemberByID(c, id)
}
