package webapi

import (
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
	return entity.Order{}, 0, nil
}
