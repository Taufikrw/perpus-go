package services

import (
	"belajar-go/dto"
	"belajar-go/models"
	"belajar-go/utils"
	"context"
)

type AuthService struct {
	tx         models.TransactionManager
	userRepo   models.UserRepository
	memberRepo models.MemberRepository
}

func NewAuthService(userRepo models.UserRepository, memberRepo models.MemberRepository, tx models.TransactionManager) *AuthService {
	return &AuthService{userRepo: userRepo, memberRepo: memberRepo, tx: tx}
}

func (s *AuthService) Login(c context.Context, input dto.LoginInput) (string, error) {
	user, err := s.userRepo.GetUserByEmail(c, input.Email)
	if err != nil {
		return "", utils.NewUnauthorizedError("Email and password do not match")
	}

	if !utils.CheckPassword(user.Password, input.Password) {
		return "", utils.NewUnauthorizedError("Email and password do not match")
	}

	token, err := utils.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Register(c context.Context, input dto.CreateMemberDTO) (*models.Member, error) {
	roleMember, err := s.userRepo.GetRoleByName(c, "member")
	if err != nil {
		return nil, utils.NewNotFoundError("Role member not found")
	}
	hashPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashPassword,
		RoleID:   roleMember.ID,
	}

	newMember := models.Member{
		MemberCode:  input.MemberCode,
		PhoneNumber: input.PhoneNumber,
		Address:     input.Address,
		IsApproved:  false,
	}

	err = s.tx.WithTransaction(c, func(txCtx context.Context) error {
		if err := s.userRepo.Create(txCtx, &user); err != nil {
			return err
		}

		newMember.UserID = user.ID
		if err := s.memberRepo.Create(txCtx, &newMember); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.memberRepo.FindByID(c, newMember.ID.String())
}
