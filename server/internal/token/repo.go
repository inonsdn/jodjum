package token

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenUser struct {
	UserId    uuid.UUID `json:"id"`
	SessionId uuid.UUID `json:"session_id"`
}

type TokenRepo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *TokenRepo {
	return &TokenRepo{
		db: db,
	}
}

func (r *TokenRepo) GetSessionOfUser(ctx context.Context, userId uuid.UUID) (TokenUser, error) {
	statement := "SELECT id, session_id FROM users WHERE id = $1"

	var tokenUser TokenUser
	err := r.db.QueryRow(ctx, statement, userId).Scan(
		&tokenUser.UserId,
		&tokenUser.SessionId,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return TokenUser{}, pgx.ErrNoRows
		}
		return TokenUser{}, err
	}

	return tokenUser, nil
}

func (r *TokenRepo) UpdateUserSession(ctx context.Context, userId uuid.UUID, sessionId uuid.UUID) (TokenUser, error) {
	statement := "UPDATE users SET session_id = $1 WHERE id = $2 RETURNING session_id"
	var tokenUser TokenUser
	err := r.db.QueryRow(ctx,
		statement,
		userId,
		sessionId).Scan(
		&tokenUser.SessionId,
	)

	if err != nil {
		return TokenUser{}, err
	}
	return tokenUser, nil
}
