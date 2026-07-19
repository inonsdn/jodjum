package auth

import (
	"context"
	"errors"
	"log/slog"
	"server/internal/constants"
	"server/internal/token"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepo     *AuthRepo
	tokenService *token.TokenService
}

type LoginResult struct {
	AccessToken string
	UserId      string
}

func NewService(AuthRepo *AuthRepo, tokenService *token.TokenService) *AuthService {
	return &AuthService{
		authRepo:     AuthRepo,
		tokenService: tokenService,
	}
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (LoginResult, error) {

	// remove space from email
	userEmail := strings.TrimSpace(email)

	// if no email input or password input, invalid credentials raise error
	if userEmail == "" || password == "" {
		return LoginResult{}, constants.ErrInvalidCredentials
	}

	// verify that email is exist by get user from email
	authUser, err := s.authRepo.GetUserFromEmail(ctx, email)
	if err != nil {
		return LoginResult{}, constants.ErrInvalidCredentials
	}

	// verify password
	if err := bcrypt.CompareHashAndPassword(authUser.Password, []byte(password)); err != nil {
		return LoginResult{}, constants.ErrInvalidCredentials
	}

	sessionId := uuid.New()

	// generate token of this users
	token, err := s.tokenService.Generate(authUser.UserId, sessionId)
	if err != nil {
		slog.Error("Cannot generate token", "Error", err)
		return LoginResult{}, constants.ErrInvalidCredentials
	}

	_, err = s.authRepo.UpdateUserSession(ctx, authUser.UserId, sessionId)
	if err != nil {
		return LoginResult{}, constants.ErrInvalidCredentials
	}

	return LoginResult{
		AccessToken: token,
		UserId:      authUser.UserId.String(),
	}, nil
}

func (s *AuthService) RegisterUser(ctx context.Context, username string, email string, password string) (LoginResult, error) {
	// remove space from email
	userEmail := strings.TrimSpace(email)

	// if no email input or password input, invalid credentials raise error
	if userEmail == "" || password == "" {
		return LoginResult{}, constants.ErrInvalidCredentials
	}

	// verify email must not duplicate in database
	authUser, err := s.authRepo.GetUserFromEmail(ctx, email)

	slog.Debug("Get user from email", "email", email)

	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("Cannot get with error", "Error", err)
			return LoginResult{}, err
		}
	}

	if authUser.UserId != uuid.Nil {
		slog.Error("Duplicated email", "Error", err)
		return LoginResult{}, constants.ErrDuplicatedEmail
	}

	// hash password
	hashPasswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// create user in database
	authUser, err = s.authRepo.CreateNewUser(ctx, username, email, hashPasswd)

	if err != nil {
		slog.Error("Cannot get with error", "Error", err)
		return LoginResult{}, err
	}
	sessionId := uuid.New()

	// generate token of this users
	token, err := s.tokenService.Generate(authUser.UserId, sessionId)
	if err != nil {
		slog.Error("Cannot generate token", "Error", err)
		return LoginResult{}, constants.ErrInvalidCredentials
	}

	return LoginResult{
		AccessToken: token,
		UserId:      authUser.UserId.String(),
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, userId uuid.UUID) error {

	sessionId := uuid.Nil
	_, err := s.authRepo.UpdateUserSession(ctx, userId, sessionId)
	return err
}
