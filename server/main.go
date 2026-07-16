package main

import (
	"fmt"
	"net/http"
	"server/internal/auth"
	"server/internal/db"
	"server/internal/user"
)

type App struct {
	server *http.Server
}

func New() *App {

	router := http.NewServeMux()

	dbConnection, err := db.NewSupabase("", "")

	if err != nil {
		fmt.Println("Found error when connect to supabase", err)
		return nil
	}

	// init user module
	userRepo := user.NewRepo(dbConnection)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)
	user.RegisterRoutes(router, userHandler)

	// init auth module
	authRepo := auth.NewRepo(dbConnection)
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
