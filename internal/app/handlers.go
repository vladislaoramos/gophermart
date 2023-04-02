package app

import (
	"encoding/json"
	"github.com/go-chi/jwtauth/v5"
	"github.com/vladislaoramos/gophemart/internal/entity"
	"github.com/vladislaoramos/gophemart/internal/usecase"
	"github.com/vladislaoramos/gophemart/pkg/logger"
	"net/http"
)

func registrationUser(
	ls usecase.LoyalSystem,
	log logger.LogInterface,
	tokenAuth *jwtauth.JWTAuth,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userAuth entity.UserAuth

		if err := json.NewDecoder(r.Body).Decode(&userAuth); err != nil {
			log.Error(err.Error())
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		user, err := ls.CreateUser(r.Context(), userAuth)
		if err != nil {
			log.Error(err.Error())
			errorHandler(w, err)
			return
		}

		_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": user.ID})
		w.Header().Set("Authorization", "Bearer "+tokenString)
		w.WriteHeader(http.StatusOK)
	}
}

func loginUser(
	ls usecase.LoyalSystem,
	log logger.LogInterface,
	tokenAuth *jwtauth.JWTAuth,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userAuth entity.UserAuth

		if err := json.NewDecoder(r.Body).Decode(&userAuth); err != nil {
			log.Error(err.Error())
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		user, err := ls.CheckUser(r.Context(), userAuth)
		if err != nil {
			errorHandler(w, err)
			return
		}

		_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": user.ID})
		w.Header().Set("Authorization", "Bearer "+tokenString)
		w.WriteHeader(http.StatusOK)
	}
}

func pingHandler(ls usecase.LoyalSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := ls.PingRepo(r.Context()); err != nil {
			http.Error(w, "repo error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func getWithdrawInfoList(_ usecase.LoyalSystem, _ logger.LogInterface, _ *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func withdraw(_ usecase.LoyalSystem, _ logger.LogInterface, _ *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func getCurrentBalance(_ usecase.LoyalSystem, _ logger.LogInterface, _ *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func getOrderInfoList(_ usecase.LoyalSystem, _ logger.LogInterface, _ *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func uploadOrder(_ usecase.LoyalSystem, _ logger.LogInterface, _ *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
