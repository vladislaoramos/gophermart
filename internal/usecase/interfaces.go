package usecase

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"github.com/vladislaoramos/gophemart/internal/entity"
)

type LoyalSystem interface {
	PingRepo(context.Context) error
	CreateUser(context.Context, entity.UserAuth) (entity.User, error)
	CheckUser(context.Context, entity.UserAuth) (entity.User, error)

	UploadOrder(context.Context, int, string) (bool, error)
	GetOrderList(context.Context, int) ([]entity.Order, error)

	GetBalance(context.Context, int) (entity.Balance, error)
	Withdraw(context.Context, int, entity.Withdrawal) error

	GetWithdrawList(context.Context, int) ([]entity.Withdraw, error)
	ProcessOrder(string) error
}

type LoyalSystemRepo interface {
	Ping(context.Context) error

	CreateUser(context.Context, string, string) (entity.User, error)
	GetUserWithLogin(context.Context, string) (entity.User, error)
	CreateUserBalance(context.Context, int) error

	GetOrderList(context.Context, int) ([]entity.Order, error)

	GetBalance(context.Context, int) (entity.Balance, error)
	UpdateBalance(context.Context, int, decimal.Decimal, decimal.Decimal) error
	AddWithdrawal(context.Context, int, string, decimal.Decimal) error

	GetWithdrawalList(context.Context, int) ([]entity.Withdraw, error)

	GetOrderByOrderNumber(context.Context, string) (entity.Order, error)
	CreateOrder(context.Context, int, string) error
	UpdateOrderStatus(context.Context, string, string) error
	UpdateOrderAccrual(context.Context, string, decimal.Decimal) error
}

type LoyalSystemAPI interface {
	GetOrderInfo(string) (entity.Order, time.Duration, error)
}
