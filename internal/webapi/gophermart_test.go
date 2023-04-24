package webapi

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
	"github.com/vladislaoramos/gophermart/internal/entity"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetOrderInfo(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/api/orders/123" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		order := entity.OrderAdapter{
			Number:  "123",
			Status:  "Delivered",
			Accrual: 10.5,
		}
		jsonResp, _ := json.Marshal(order)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	}))
	defer mockServer.Close()

	loyalSystemWebAPI := LoyalSystemWebAPI{
		client: resty.New().SetBaseURL(mockServer.URL),
	}

	order, _, err := loyalSystemWebAPI.GetOrderInfo("123")
	require.NoError(t, err)

	expectedOrder := entity.Order{
		Number:  "123",
		Status:  "Delivered",
		Accrual: 10.5,
	}

	require.EqualValues(t, expectedOrder, order)
}

func TestGetOrderInfo_ManyRequests(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/api/orders/123" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer mockServer.Close()

	loyalSystemWebAPI := LoyalSystemWebAPI{
		client: resty.New().SetBaseURL(mockServer.URL),
	}

	_, _, err := loyalSystemWebAPI.GetOrderInfo("123")

	require.Error(t, err, &ErrTooManyRequests)
}

func TestGetOrderInfo_InternalError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/api/orders/123" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	loyalSystemWebAPI := LoyalSystemWebAPI{
		client: resty.New().SetBaseURL(mockServer.URL),
	}

	_, _, err := loyalSystemWebAPI.GetOrderInfo("123")

	require.Error(t, err, &ErrInternalServerError)
}
