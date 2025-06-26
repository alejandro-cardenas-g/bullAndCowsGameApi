package api

import (
	"log"
	"net/http"

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
	Addr string
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

func (app *Application) Run() error {
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

	log.Println("Listening on", app.config.Addr)

	return http.ListenAndServe(app.config.Addr, router)
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
