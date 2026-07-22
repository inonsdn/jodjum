package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"server/internal/auth"
	"server/internal/config"
	"server/internal/db"
	"server/internal/middleware"
	"server/internal/notification"
	"server/internal/reminders"
	"server/internal/response"
	"server/internal/things"
	"server/internal/token"
	"server/internal/user"
	"syscall"
	"time"
)

type App struct {
	server      *http.Server
	reminderApp *reminders.ReminderApp
	notifier    notification.Notification
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

	slog.Info("Init reminders module")
	// init reminders module (before things: things creates reminders in the
	// same transaction, so it needs the reminders repo)
	remindersRepo := reminders.NewRepo(con)
	remindersService := reminders.NewService(remindersRepo)
	remindersHandler := reminders.NewHandler(remindersService)
	reminders.RegisterRoutes(router, authService.AuthMiddleware, remindersHandler)

	slog.Info("Init things module")
	// init things module
	thingsRepo := things.NewRepo(con)
	thingsService := things.NewService(thingsRepo, remindersRepo)
	thingsHandler := things.NewHandler(thingsService)
	things.RegisterRoutes(router, authService.AuthMiddleware, thingsHandler)

	slog.Info("Init notification module")
	// init notification module (Web Push subscription management)
	notificationRepo := notification.NewRepo(con)
	notificationService := notification.NewService(notificationRepo)
	notificationHandler := notification.NewHandler(notificationService)
	notification.RegisterRoutes(router, authService.AuthMiddleware, notificationHandler)

	slog.Info("Init reminder loop")
	// Web Push sender + background reminder app (goroutine started in App.Run).
	webPush := notification.NewWebPushNotification(*cfg.GetWebPushNotificationConfig(), notificationService)
	reminderApp := reminders.NewReminderApp(remindersRepo)

	server := http.Server{
		Addr: cfg.Address(),
		// Logger is outermost so it records every request (including the CORS
		// preflight OPTIONS that CORS answers). CORS then wraps the router so
		// it runs before auth.
		Handler: middleware.Logger(middleware.CORS(cfg.AllowedOrigin)(router)),
	}

	return &App{
		server:      &server,
		reminderApp: reminderApp,
		notifier:    webPush,
	}
}

func (a *App) Run(ctx context.Context) {
	// Start the background reminder loop; it stops when ctx is cancelled.
	a.reminderApp.Run(ctx, a.notifier)

	// Serve in a goroutine so we can block on the shutdown signal below.
	go func() {
		slog.Info("Run and serve", "address", a.server.Addr)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
		}
	}()

	// Block until a stop signal cancels ctx.
	<-ctx.Done()
	slog.Info("shutting down")

	// Give in-flight requests up to 10s to finish.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.server.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	}
}

func main() {
	initLogger()

	// ctx is cancelled when the OS sends SIGINT (Ctrl-C) or SIGTERM (what
	// Cloud Run / Fly send to stop the container).
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app := New()
	if app == nil {
		return
	}
	app.Run(ctx)
}
