package auth

import (
	"context"
	"errors"
	"net/http"
	"server/internal/response"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *AuthService) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get authorization from header
		header := r.Header.Get("Authorization")

		// extract Bearer scheme from authorization
		tokenString, isFound := strings.CutPrefix(header, "Bearer ")
		tokenString = strings.TrimSpace(tokenString)
		if !isFound {
			response.JSON(w, http.StatusUnauthorized, "missing or invalid authorization header")
			return
		}

		// get token and extract user info from token
		parseResult, err := s.tokenService.Parse(tokenString)
		if err != nil {
			response.JSON(w, http.StatusUnauthorized, err.Error())
			return
		}

		// check session id is valid
		authUser, err := s.authRepo.GetSessionOfUser(r.Context(), parseResult.UserId)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				response.JSON(w, http.StatusUnauthorized, "No user found, create new user first")
			} else {
				response.JSON(w, http.StatusUnauthorized, err.Error())
			}
			return
		}

		if authUser.SessionId == uuid.Nil {
			response.JSON(w, http.StatusUnauthorized, "No login session, please login first")
			return
		}

		if parseResult.SessionId != authUser.SessionId {
			response.JSON(w, http.StatusUnauthorized, "Session is no longer used, please login first")
			return
		}

		ctx := context.WithValue(
			r.Context(),
			UserIDContextKey,
			parseResult.UserId,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
