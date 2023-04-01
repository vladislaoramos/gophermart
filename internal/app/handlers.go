package app

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/vladislaoramos/gophemart/internal/usecase"
	"github.com/vladislaoramos/gophemart/pkg/logger"
	"net/http"
)

func getWithdrawInfoList(_ usecase.Gophermart, _ logger.LogInterface, _ *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func withdraw(_ usecase.Gophermart, _ logger.LogInterface, _ *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func getCurrentBalance(_ usecase.Gophermart, _ logger.LogInterface, _ *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func getOrderInfoList(_ usecase.Gophermart, _ logger.LogInterface, _ *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func uploadOrder(_ usecase.Gophermart, _ logger.LogInterface, _ *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func registrationUser(_ usecase.Gophermart, _ logger.LogInterface, _ *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func loginUser(_ usecase.Gophermart, _ logger.LogInterface, _ *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func pingHandler(_ usecase.Gophermart, _ logger.LogInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
