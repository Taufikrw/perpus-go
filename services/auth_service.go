package services

import (
	"belajar-go/dto"
	"belajar-go/models"
	"belajar-go/utils"
	"context"
)

type AuthService struct {
	repository models.AuthRepositoryInterface
}

func NewAuthService(repository models.AuthRepositoryInterface) *AuthService {
	return &AuthService{repository: repository}
}

func (s *AuthService) Login(c context.Context, input dto.LoginInput) (string, error) {
	user, err := s.repository.GetUserByEmail(c, input.Email)
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
	roleMember, err := s.repository.GetRoleByName(c, "member")
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

	return s.repository.RegisterMemberTransaction(c, &user, &newMember)
}
