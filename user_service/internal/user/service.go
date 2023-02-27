package user

import (
	"context"
	"fmt"
	"user_service/internal/domain"

	"github.com/Levap123/utils/jwt"
)

type UserService struct {
	repo IUserRepo
	j    *jwt.JWT
}

func NewUserService(repo IUserRepo, j *jwt.JWT) *UserService {
	return &UserService{
		repo: repo,
		j:    j,
	}
}

type IUserRepo interface {
	Create(ctx context.Context, user *User) (uint64, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

func (us *UserService) Create(ctx context.Context, dto *CreateUserDTO) (uint64, error) {
	user := NewUserFromCreateDTO(dto)
	if err := user.generatePasswordHash(); err != nil {
		return 0, fmt.Errorf("user service - generate hash - %w", err)
	}
	userID, err := us.repo.Create(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("user service - repo create - %w", err)
	}
	return userID, nil
}

const (
	accessType  = "access"
	refreshType = "refresh"
)

func (us *UserService) GenerateTokens(ctx context.Context, dto *GetUserDTO) (string, string, error) {
	user, err := us.repo.GetByEmail(ctx, dto.Email)
	if err != nil {
		return "", "", fmt.Errorf("user service - %w", err)
	}

	if !user.PasswordCorrect(dto.Password) {
		return "", "", domain.ErrIncorrectPassword
	}

	accessToken, err := us.j.GenerateJwt(int(user.ID), 2, accessType)
	if err != nil {
		return "", "", fmt.Errorf("user service - generate access token - %w", err)
	}

	refreshToken, err := us.j.GenerateJwt(int(user.ID), 2, refreshType)
	if err != nil {
		return "", "", fmt.Errorf("user service - generate refresh token - %w", err)
	}

	return accessToken, refreshToken, nil
}

func (us *UserService) Validate(ctx context.Context, accessToken string) (int, error) {
	claims, err := us.j.ParseToken(accessToken)
	if err != nil {
		return 0, err
	}
	if claims.TokenType != accessType {
		return 0, domain.ErrIncorrectTokenType
	}
	return claims.UserID, nil
}
