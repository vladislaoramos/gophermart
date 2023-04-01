package webapi

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/vladislaoramos/gophemart/internal/entity"
)

type GophermartWebAPI struct {
	client *resty.Client
}

func New(client *resty.Client) *GophermartWebAPI {
	return &GophermartWebAPI{
		client: client,
	}
}

func (w *GophermartWebAPI) GetOrderInfo(orderNumber string) (entity.Order, time.Duration, error) {
	return entity.Order{}, 0, nil
}
