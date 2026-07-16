package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthService struct {
	authRepo *AuthRepo
}

type LoginResult struct {
	AccessToken string
	UserId      string
}

func NewService(AuthRepo *AuthRepo) *AuthService {
	return &AuthService{
		authRepo: AuthRepo,
	}
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (LoginResult, error) {

	// remove space from email
	userEmail := strings.TrimSpace(email)

	_, err := s.authRepo.db.Login(email, password)

	if err != nil {
		return LoginResult{}, ErrInvalidCredentials
	}

	// if no email input or password input, invalid credentials raise error
	if userEmail == "" || password == "" {
		return LoginResult{}, ErrInvalidCredentials
	}

	// verify that email is exist by get user from email

	// verify password
	if err := bcrypt.CompareHashAndPassword([]byte("abc"), []byte(password)); err != nil {
		return LoginResult{}, ErrInvalidCredentials
	}

	// generate token of this users

	return LoginResult{
		AccessToken: "test",
		UserId:      "1234124a",
	}, nil
}

func (s *AuthService) RegisterUser(ctx context.Context, username string, email string, password string) (LoginResult, error) {
	// remove space from email
	userEmail := strings.TrimSpace(email)

	// if no email input or password input, invalid credentials raise error
	if userEmail == "" || password == "" {
		return LoginResult{}, ErrInvalidCredentials
	}

	// verify email must not duplicate in database

	// hash password
	hashPasswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// create user in database

	return LoginResult{}, nil
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}
type TokenService struct {
	secret []byte
}

func (s *TokenService) Generate() (string, error) {
	return "", nil
}
