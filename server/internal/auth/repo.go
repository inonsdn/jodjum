package auth

import "server/internal/db"

type AuthRepo struct {
	db db.DbConnection
}

func NewRepo(db db.DbConnection) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}
