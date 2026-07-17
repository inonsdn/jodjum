package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrDuplicatedEmail    = errors.New("email is already used")
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

	// if no email input or password input, invalid credentials raise error
	if userEmail == "" || password == "" {
		return LoginResult{}, ErrInvalidCredentials
	}

	// verify that email is exist by get user from email
	authUser, err := s.authRepo.GetUserFromEmail(ctx, email)
	if err != nil {
		return LoginResult{}, ErrInvalidCredentials
	}

	// verify password
	if err := bcrypt.CompareHashAndPassword(authUser.Password, []byte(password)); err != nil {
		return LoginResult{}, ErrInvalidCredentials
	}

	// generate token of this users

	return LoginResult{
		AccessToken: "test",
		UserId:      authUser.UserId.String(),
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
	authUser, err := s.authRepo.GetUserFromEmail(ctx, email)

	fmt.Println("Get user from email", email)

	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			fmt.Println("Cannot get with error", err)
			return LoginResult{}, err
		}
	}

	if authUser.UserId != uuid.Nil {
		fmt.Println("Duplicated email", err)
		return LoginResult{}, ErrDuplicatedEmail
	}

	// hash password
	hashPasswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// create user in database
	s.authRepo.CreateNewUser(ctx, username, email, hashPasswd)

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
