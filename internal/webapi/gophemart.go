package webapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/vladislaoramos/gophemart/internal/entity"
)

type LoyalSystemWebAPI struct {
	client *resty.Client
}

func NewAPI(client *resty.Client) *LoyalSystemWebAPI {
	return &LoyalSystemWebAPI{
		client: client,
	}
}

func (w *LoyalSystemWebAPI) GetOrderInfo(orderNumber string) (entity.Order, time.Duration, error) {
	var order entity.OrderAdapter

	resp, err := w.client.
		R().
		SetResult(&order).
		Get("/api/orders/" + orderNumber)
	if err != nil {
		return entity.Order{}, 0, fmt.Errorf("WebAPI - GetOrderInfo - w.client.R().Get: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusTooManyRequests:
		timeout, err := time.ParseDuration(resp.Header().Get("Retry-After") + "s")
		if err != nil {
			return entity.Order{}, 0, fmt.Errorf("WebAPI - GetOrderInfo - time.ParseDuration: %w", err)
		}
		return entity.Order{}, timeout, ErrTooManyRequests
	case http.StatusInternalServerError:
		return entity.Order{}, 0, ErrInternalServerError
	default:
		return entity.Order{
			Number:  order.Number,
			Status:  order.Status,
			Accrual: order.Accrual,
		}, 0, nil
	}
}
