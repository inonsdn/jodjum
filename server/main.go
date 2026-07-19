package main

import (
	"log/slog"
	"net/http"
	"os"
	"server/internal/auth"
	"server/internal/config"
	"server/internal/db"
	"server/internal/middleware"
	"server/internal/reminders"
	"server/internal/response"
	"server/internal/things"
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

	// Public health check for load balancers / Cloud Run probes (no auth).
	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("health check hit")
		response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

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

	slog.Info("Init things module")
	// init things module
	thingsRepo := things.NewRepo(con)
	thingsService := things.NewService(thingsRepo)
	thingsHandler := things.NewHandler(thingsService)
	things.RegisterRoutes(router, authService.AuthMiddleware, thingsHandler)

	slog.Info("Init reminders module")
	// init reminders module
	remindersRepo := reminders.NewRepo(con)
	remindersService := reminders.NewService(remindersRepo)
	remindersHandler := reminders.NewHandler(remindersService)
	reminders.RegisterRoutes(router, authService.AuthMiddleware, remindersHandler)

	server := http.Server{
		Addr: cfg.Address(),
		// Logger is outermost so it records every request (including the CORS
		// preflight OPTIONS that CORS answers). CORS then wraps the router so
		// it runs before auth.
		Handler: middleware.Logger(middleware.CORS(cfg.AllowedOrigin)(router)),
	}

	return &App{
		server: &server,
	}
}

func (a *App) Run() {
	slog.Info("Run and serve", "address", a.server.Addr)
	if err := a.server.ListenAndServe(); err != nil {
		slog.Error("ERROR", "error", err)
	}
}

func main() {
	initLogger()
	app := New()
	if app == nil {
		return
	}
	app.Run()
}
