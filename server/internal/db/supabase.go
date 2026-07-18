package db

import (
	"log/slog"

	"github.com/supabase-community/supabase-go"
)

type Supabase struct {
	client *supabase.Client
}

func NewSupabase(url string, key string) (*Supabase, error) {

	opt := supabase.ClientOptions{}
	client, err := supabase.NewClient(url, key, &opt)

	if err != nil {
		return nil, err
	}

	return &Supabase{
		client: client,
	}, nil
}

func (s *Supabase) Login(email string, password string) (AuthTokenSession, error) {
	session, err := s.client.SignInWithEmailPassword(email, password)

	slog.Info("LOGIN ", session)

	if err != nil {
		return AuthTokenSession{}, err
	}

	return AuthTokenSession{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		TokenType:    session.TokenType,
		ExpiresIn:    session.ExpiresIn,
		ExpiresAt:    session.ExpiresAt,
		UserId:       session.User.ID,
	}, nil
}

func (s *Supabase) Query() {
	data, count, _ := s.client.From("user").Select("*", "exact", false).Execute()

	slog.Info("DATA", data)
	slog.Info("COUNT", count)
}
