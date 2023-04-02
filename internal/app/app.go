package app

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/vladislaoramos/gophemart/configs"
	"github.com/vladislaoramos/gophemart/internal/repo"
	"github.com/vladislaoramos/gophemart/internal/usecase"
	"github.com/vladislaoramos/gophemart/internal/webapi"
	"github.com/vladislaoramos/gophemart/pkg/logger"
	"github.com/vladislaoramos/gophemart/pkg/postgres"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const workersCount = 5

func Run(cfg *configs.Config, log *logger.Logger) {
	log.Info("App is running...")
	err := applyMigration(cfg.Database.URI)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - migration: %w", err).Error())
	}

	db, err := postgres.New(cfg.Database.URI)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - postgres.NewAPI: %w", err).Error())
	}
	defer db.Close()

	client := resty.New().SetBaseURL(cfg.App.AccrualSystemAddress)
	webAPI := webapi.NewAPI(client)
	repository := repo.NewRepository(db)

	uc := usecase.NewLoyalSystem(
		repository,
		webAPI,
		workersCount,
		log,
	)

	handler := chi.NewRouter()
	NewRouter(handler, uc, log)

	log.Fatal(http.ListenAndServe(cfg.Client.Address, handler).Error())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	stop := <-sigs

	log.Info("App got stop signal: " + stop.String())
}
