package services

import (
	"errors"
	"strconv"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/utils"
)

type User struct {
	ID      uint
	Name    string
	Email   string
	IsSuper bool
}

type Response struct {
	AccessToken  string
	RefreshToken string
	User         User
}

type AuthService struct {
	UserService *UserService
}

func NewAuthService(userService *UserService) *AuthService {
	return &AuthService{UserService: userService}
}

func (s *AuthService) Login(email string, password string) (Response, error) {
	if email == "" {
		return Response{}, errors.New("EMAIL_IS_REQUIRED")
	}

	if password == "" {
		return Response{}, errors.New("PASSWORD_IS_REQUIRED")
	}

	user, err := s.UserService.FindByEmail(email)

	if err != nil {
		return Response{}, err
	}

	if err := user.VerifyPassword(password); err != nil {
		return Response{}, errors.New("INVALID_PASSWORD")
	}

	var accessToken string
	if accessToken, err = utils.CreateAccessToken(user.ID); err != nil {
		return Response{}, err
	}

	var refreshToken string
	if refreshToken, err = utils.CreateRefreshToken(user.ID); err != nil {
		return Response{}, err
	}

	return Response{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: User{
			ID:      user.ID,
			Name:    user.Name,
			Email:   user.Email,
			IsSuper: user.IsSuper,
		},
	}, nil
}

func (s *AuthService) Refresh(token string) (Response, error) {
	parsed, err := utils.VerifyToken(token)

	if err != nil {
		return Response{}, err
	}

	sub, err := parsed.Claims.GetSubject()

	if err != nil {
		return Response{}, err
	}

	id, err := strconv.ParseUint(sub, 10, 32)
	if err != nil {
		return Response{}, errors.New("INVALID_SUBJECT")
	}

	user, err := s.UserService.FindById(uint(id))

	if err != nil {
		return Response{}, err
	}

	var accessToken string
	if accessToken, err = utils.CreateAccessToken(user.ID); err != nil {
		return Response{}, err
	}

	var refreshToken string
	if refreshToken, err = utils.CreateRefreshToken(user.ID); err != nil {
		return Response{}, err
	}

	return Response{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: User{
			ID:      user.ID,
			Name:    user.Name,
			Email:   user.Email,
			IsSuper: user.IsSuper,
		},
	}, nil
}
