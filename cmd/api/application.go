package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/alejandro-cardenas-g/bullAndCowsApp/internal/services"
	"github.com/alejandro-cardenas-g/bullAndCowsApp/internal/store"
	"github.com/alejandro-cardenas-g/bullAndCowsApp/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

type ApplicationConfig struct {
	Addr            string
	GracefulTimeout time.Duration
}

type Application struct {
	config ApplicationConfig
	logger *zap.SugaredLogger
}

func NewApplication(
	config ApplicationConfig,
) *Application {

	logger := zap.Must(zap.NewProduction()).Sugar()

	return &Application{
		config: config,
		logger: logger,
	}
}

func (app *Application) Run() {

	router := app.createRouter()

	srv := &http.Server{
		Addr:         app.config.Addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
		app.logger.Info("Listening on", app.config.Addr)
	}()

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt)

	<-ch

	ctx, cancel := context.WithTimeout(context.Background(), app.config.GracefulTimeout)
	defer cancel()

	srv.Shutdown(ctx)
	os.Exit(0)
}

func (app *Application) createRouter() *mux.Router {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	controller := &Controller{
		logger: app.logger,
	}

	// storage registration
	matchesRdb := app.createMatchesRdb()
	storage := store.NewRedisStorage(matchesRdb)

	// services registration
	matchesService := services.NewMatchesService(storage)

	// controllers registration
	matchesController := newMatchesController(controller, matchesService)
	matchesController.RegisterRoutes(subrouter)
	return router
}

func (app *Application) createMatchesRdb() *redis.Client {
	matchesRdb, err := store.NewRedisClient(
		utils.GetEnvironment().GetEnv("DB_MATCHES", ""),
		utils.GetEnvironment().GetEnv("DB_MATCHES_PWD", ""),
		int(utils.GetEnvironment().GetEnvAsInt("DB_MATCHES_DB", 0)),
	)
	if err != nil {
		app.logger.Error("error connecting to matchesdb")
		log.Fatal(err)
	}
	return matchesRdb
}
