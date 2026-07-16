package user

import "server/internal/db"

type UserRepo struct {
	db db.DbConnection
}

func NewRepo(db db.DbConnection) *UserRepo {
	return &UserRepo{
		db: db,
	}
}
