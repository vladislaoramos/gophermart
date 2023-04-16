package app

import (
	"encoding/json"
	"github.com/go-chi/jwtauth/v5"
	"github.com/theplant/luhn"
	"github.com/vladislaoramos/gophermart/internal/entity"
	"github.com/vladislaoramos/gophermart/internal/usecase"
	"github.com/vladislaoramos/gophermart/pkg/logger"
	"io"
	"net/http"
	"strconv"
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

func getWithdrawList(ls usecase.LoyalSystem, log logger.LogInterface, _ *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		wl, err := ls.GetWithdrawList(r.Context(), int(claims["user_id"].(float64)))
		if err != nil {
			log.Error(err.Error())
			errorHandler(w, err)
			return
		}

		jsonResp, err := json.Marshal(wl)
		if err != nil {
			log.Error(err.Error())
			errorHandler(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	}
}

func withdraw(ls usecase.LoyalSystem, log logger.LogInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var withdrawal entity.Withdrawal

		if err := json.NewDecoder(r.Body).Decode(&withdrawal); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		orderNum, _ := strconv.Atoi(withdrawal.Order)
		if !luhn.Valid(orderNum) {
			http.Error(w, "bad request", http.StatusUnprocessableEntity)
			return
		}

		_, claims, _ := jwtauth.FromContext(r.Context())
		err := ls.Withdraw(r.Context(), int(claims["user_id"].(float64)), withdrawal)
		if err != nil {
			log.Error(err.Error())
			errorHandler(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func getBalance(ls usecase.LoyalSystem, log logger.LogInterface, _ *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		balance, err := ls.GetBalance(r.Context(), int(claims["user_id"].(float64)))
		if err != nil {
			log.Error(err.Error())
			errorHandler(w, err)
			return
		}

		jsonResp, err := json.Marshal(balance)
		if err != nil {
			log.Error(err.Error())
			errorHandler(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	}
}

func getOrderInfoList(ls usecase.LoyalSystem, log logger.LogInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		orderList, err := ls.GetOrderList(r.Context(), int(claims["user_id"].(float64)))
		if err != nil {
			log.Error(err.Error())
			errorHandler(w, err)
			return
		}

		if orderList == nil {
			http.Error(w, "empty order list", http.StatusNoContent)
			return
		}

		jsonResp, err := json.Marshal(orderList)
		if err != nil {
			log.Error(err.Error())
			errorHandler(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	}
}

func uploadOrder(
	ls usecase.LoyalSystem,
	log logger.LogInterface,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		orderNum, _ := strconv.Atoi(string(body))
		if !luhn.Valid(orderNum) {
			http.Error(w, "bad request", http.StatusUnprocessableEntity)
			return
		}

		_, claims, _ := jwtauth.FromContext(r.Context())
		isAlreadyUploaded, err := ls.UploadOrder(r.Context(), int(claims["user_id"].(float64)), strconv.Itoa(orderNum))
		if err != nil {
			errorHandler(w, err)
			return
		}

		if isAlreadyUploaded {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusAccepted)
		}
	}
}
