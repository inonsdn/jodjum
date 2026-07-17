package main

import (
	"fmt"
	"net/http"
	"os"
	"server/internal/auth"
	"server/internal/db"
	"server/internal/user"

	"github.com/joho/godotenv"
)

type App struct {
	server *http.Server
}

func New() *App {

	if err := godotenv.Load(); err != nil {
		fmt.Println(".env file not found, using system environment")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	con, err := db.NewPGX(databaseURL)

	router := http.NewServeMux()

	if err != nil {
		fmt.Println("Found error when connect to supabase", err)
		return nil
	}

	// init user module
	userRepo := user.NewRepo(con)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)
	user.RegisterRoutes(router, userHandler)

	// init auth module
	authRepo := auth.NewRepo(con)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)
	auth.RegisterRoutes(router, authHandler)

	server := http.Server{
		Addr:    "localhost:8000",
		Handler: router,
	}

	return &App{
		server: &server,
	}
}

func (a *App) Run() {
	a.server.ListenAndServe()
}

func main() {
	app := New()
	if app == nil {
		return
	}
	app.Run()
}
