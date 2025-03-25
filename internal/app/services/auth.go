package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func (s *AuthService) CreateAccessToken(id uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

func (s *AuthService) CreateRefreshToken(id uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

func (s *AuthService) VerifyToken(token string) error {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return os.Getenv("SECRET_KEY"), nil
	})

	if err != nil {
		return err
	}

	if !parsed.Valid {
		return fmt.Errorf("INVALID_TOKEN")
	}

	return nil
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
	if accessToken, err = s.CreateAccessToken(user.ID); err != nil {
		return Response{}, err
	}

	var refreshToken string
	if refreshToken, err = s.CreateRefreshToken(user.ID); err != nil {
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
