package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthUser struct {
	UserId    uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  []byte    `json:"password"`
	SessionId uuid.UUID `json:"session_id"`
}

type AuthRepo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (r *AuthRepo) GetUserFromEmail(ctx context.Context, email string) (AuthUser, error) {
	statement := "SELECT id, email, password FROM users WHERE email = $1"

	var authUser AuthUser
	err := r.db.QueryRow(ctx, statement, email).Scan(
		&authUser.UserId,
		&authUser.Email,
		&authUser.Password,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AuthUser{}, pgx.ErrNoRows
		}
		return AuthUser{}, err
	}

	return authUser, nil
}

func (r *AuthRepo) CreateNewUser(ctx context.Context, username string, email string, password []byte) (AuthUser, error) {
	statement := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, email"
	var authUser AuthUser
	err := r.db.QueryRow(ctx,
		statement,
		username,
		email,
		password).Scan(
		&authUser.UserId,
		&authUser.Email,
	)

	if err != nil {
		return AuthUser{}, err
	}
	return authUser, nil
}

func (r *AuthRepo) GetSessionOfUser(ctx context.Context, userId uuid.UUID) (AuthUser, error) {
	statement := "SELECT id, session_id FROM users WHERE id = $1"

	var authUser AuthUser
	err := r.db.QueryRow(ctx, statement, userId).Scan(
		&authUser.UserId,
		&authUser.SessionId,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AuthUser{}, pgx.ErrNoRows
		}
		return AuthUser{}, err
	}

	return authUser, nil
}

func (r *AuthRepo) UpdateUserSession(ctx context.Context, userId uuid.UUID, sessionId uuid.UUID) (AuthUser, error) {
	statement := "UPDATE users SET session_id = $1 WHERE id = $2 RETURNING id, session_id"
	var authUser AuthUser
	err := r.db.QueryRow(ctx,
		statement,
		sessionId,
		userId).Scan(
		&authUser.UserId,
		&authUser.SessionId,
	)

	if err != nil {
		return AuthUser{}, err
	}
	return authUser, nil
}
