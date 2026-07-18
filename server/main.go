package main

import (
	"log/slog"
	"net/http"
	"os"
	"server/internal/auth"
	"server/internal/config"
	"server/internal/db"
	"server/internal/token"
	"server/internal/user"
)

type App struct {
	server *http.Server
}

func initLogger() {
	logger := slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	))
	slog.SetDefault(logger)
}

func New() *App {

	slog.Info("Loading config")
	cfg := config.LoadConfig()
	if cfg == nil {
		return nil
	}

	slog.Info("Connect to database")
	con, err := db.NewPGX(cfg.DatabaseUrl)

	router := http.NewServeMux()

	if err != nil {
		slog.Error("Found error when connect to database", "Error", err.Error())
		return nil
	}

	slog.Info("Init auth module")
	// init auth module

	tokenRepo := token.NewRepo(con)
	tokenService := token.NewService(cfg.GetTokenConfig(), *tokenRepo)

	authRepo := auth.NewRepo(con)
	authService := auth.NewService(authRepo, tokenService)
	authHandler := auth.NewHandler(authService)
	auth.RegisterRoutes(router, authService.AuthMiddleware, authHandler)

	slog.Info("Init user module")
	// init user module
	userRepo := user.NewRepo(con)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)
	user.RegisterRoutes(router, authService.AuthMiddleware, userHandler)

	server := http.Server{
		Addr:    cfg.Address(),
		Handler: router,
	}

	return &App{
		server: &server,
	}
}

func (a *App) Run() {
	slog.Info("Run and serve", "address", a.server.Addr)
	a.server.ListenAndServe()
}

func main() {
	initLogger()
	app := New()
	if app == nil {
		return
	}
	app.Run()
}
