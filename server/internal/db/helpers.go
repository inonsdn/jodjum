package db

import "github.com/google/uuid"

type AuthTokenSession struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int
	ExpiresAt    int64
	UserId       uuid.UUID
}

type DbConnection interface {
	Query()
	Login(string, string) (AuthTokenSession, error)
}
