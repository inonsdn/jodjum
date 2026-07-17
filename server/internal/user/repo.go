package user

import "github.com/jackc/pgx/v5/pgxpool"

type UserRepo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
	}
}
