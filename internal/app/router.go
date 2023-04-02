package app

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/vladislaoramos/gophemart/internal/usecase"
	"github.com/vladislaoramos/gophemart/pkg/logger"
)

func NewRouter(handler *chi.Mux, ls usecase.LoyalSystem, log logger.LogInterface) {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	handler.Use(middleware.Logger)

	// checker
	handler.Get("/ping", pingHandler(ls))

	// auth
	handler.Group(func(r chi.Router) {
		r.Post("/api/user/register", registrationUser(ls, log, tokenAuth))
		r.Post("/api/user/login", loginUser(ls, log, tokenAuth))
	})

	// Protected routes
	handler.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Post("/api/user/orders", uploadOrder(ls, log, tokenAuth))
		r.Get("/api/user/orders", getOrderInfoList(ls, log, tokenAuth))

		r.Get("/api/user/balance", getCurrentBalance(ls, log, tokenAuth))
		r.Post("/api/user/balance/withdraw", withdraw(ls, log, tokenAuth))

		r.Get("/api/user/withdrawals", getWithdrawInfoList(ls, log, tokenAuth))
	})
}
