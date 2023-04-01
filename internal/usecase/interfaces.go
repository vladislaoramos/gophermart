package usecase

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"github.com/vladislaoramos/gophemart/internal/entity"
)

type Gophermart interface {
	PingRepo(context.Context) error
	CreateNewUser(context.Context, entity.UserAuth) (entity.User, error)
	CheckUser(context.Context, entity.UserAuth) (entity.User, error)

	UploadOrder(context.Context, int, string) (bool, error)
	GetOrderList(context.Context, int) ([]entity.Order, error)

	GetCurrentBalance(context.Context, int) (entity.Balance, error)
	Withdraw(context.Context, int, entity.Withdrawal) error

	GetWithdrawList(context.Context, int) ([]entity.Withdraw, error)
	ProcessOrder(string) error
}

type GophermartRepo interface {
	Ping(context.Context) error

	CreateUser(context.Context, string, string) (entity.User, error)
	GetUserWithLogin(context.Context, string) (entity.User, error)
	CreateUserBalance(context.Context, int) error

	GetOrderList(context.Context, int) ([]entity.Order, error)

	GetCurrentBalance(context.Context, int) (entity.Balance, error)
	UpdateBalance(context.Context, int, decimal.Decimal, decimal.Decimal) error
	AddWithdrawal(context.Context, int, string, decimal.Decimal) error

	GetWithdrawalList(context.Context, int, string) ([]entity.Withdraw, error)

	GetOrderByOrderNumber(context.Context, string) (entity.Order, error)
	CreateOrder(context.Context, int, string) error
	UpdateOrderStatus(context.Context, string, string) error
	UpdateOrderAccrual(context.Context, string, decimal.Decimal) error
}

type GophermartWebAPI interface {
	GetOrderInfo(string) (entity.Order, time.Duration, error)
}
