package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"test-case/internal/config"
	"test-case/internal/server/router"
	"test-case/storage/postgres"
	"test-case/storage/repos"
	"time"

	"github.com/go-chi/chi/v5"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

type App struct {
	Cfg      config.Config
	Storage  *postgres.Database
	UserRepo repos.UserRepository
	Router   *chi.Mux
	Server   *http.Server
}

func (app *App) readConfig() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Usage go run <path to main.go> [arguments] \n Required arguments: \n - Path to config file")
		os.Exit(1)
	}

	app.Cfg = config.ReadConfig(args[0])
}

func (app *App) SetConfig() {
	app.readConfig()

	storage, err := postgres.New(app.Cfg)
	if err != nil {

		os.Exit(1)
	}
	app.Storage = storage

	app.UserRepo = repos.NewUserRepository(app.Storage.Database)

	app.Router = router.NewRouter(app.UserRepo)

	app.Server = &http.Server{
		Addr:        app.Cfg.Address,
		Handler:     app.Router,
		ReadTimeout: app.Cfg.Timeout,
		IdleTimeout: app.Cfg.IdleTimeout,
	}
}

func (app *App) Run() {
	const op = "app.Run"

	if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalln("Fatal error", op, err.Error())
		return
	}
}

func (app *App) Stop() {
	const op = "app.Stop"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.Storage.Stop()

	if err := app.Server.Shutdown(ctx); err != nil {
		log.Println("Server forced to shutdown", op, err.Error())
		return
	}

	log.Println("Server stopped")
}

func New() *App {
	return &App{}
}
