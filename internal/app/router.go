package app

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/vladislaoramos/gophemart/internal/usecase"
	"github.com/vladislaoramos/gophemart/pkg/logger"
)

func NewRouter(handler *chi.Mux, uc usecase.Gophermart, l logger.LogInterface) {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	handler.Use(middleware.Logger)

	// checker
	handler.Get("/ping", pingHandler(uc, l))

	// auth
	handler.Group(func(r chi.Router) {
		r.Post("/api/user/register", registrationUser(uc, l, tokenAuth))
		r.Post("/api/user/login", loginUser(uc, l, tokenAuth))
	})

	// Protected routes
	handler.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Post("/api/user/orders", uploadOrder(uc, l, tokenAuth))
		r.Get("/api/user/orders", getOrderInfoList(uc, l, tokenAuth))

		r.Get("/api/user/balance", getCurrentBalance(uc, l, tokenAuth))
		r.Post("/api/user/balance/withdraw", withdraw(uc, l, tokenAuth))

		r.Get("/api/user/withdrawals", getWithdrawInfoList(uc, l, tokenAuth))
	})
}
