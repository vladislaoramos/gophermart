package repo

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/vladislaoramos/gophemart/internal/entity"
	"github.com/vladislaoramos/gophemart/pkg/postgres"
)

type GophermartRepo struct {
	*postgres.DB
}

func New(db *postgres.DB) *GophermartRepo {
	return &GophermartRepo{db}
}

func (r *GophermartRepo) Ping(_ context.Context) error {
	return nil
}

func (r *GophermartRepo) UpdateOrderAccrual(
	ctx context.Context, orderNumber string, accrual decimal.Decimal) error {
	return nil
}

func (r *GophermartRepo) UpdateOrderStatus(
	ctx context.Context, orderNumber, status string) error {
	return nil
}

func (r *GophermartRepo) GetOrderByOrderNumber(
	ctx context.Context, orderNumber string) (entity.Order, error) {
	return entity.Order{}, nil
}

func (r *GophermartRepo) CreateUserBalance(
	ctx context.Context, userID int) error {
	return nil
}

func (r *GophermartRepo) CreateOrder(
	ctx context.Context, userID int, orderNumber string) error {
	return nil
}

func (r *GophermartRepo) CreateUser(
	ctx context.Context, login, passwordHash string) (entity.User, error) {
	return entity.User{}, nil
}

func (r *GophermartRepo) GetUserWithLogin(
	ctx context.Context, login string) (entity.User, error) {
	return entity.User{}, nil
}

func (r *GophermartRepo) GetOrderList(
	ctx context.Context, userID int) ([]entity.Order, error) {
	return nil, nil
}

func (r *GophermartRepo) GetCurrentBalance(
	ctx context.Context, userID int) (entity.Balance, error) {
	return entity.Balance{}, nil
}

func (r *GophermartRepo) UpdateBalance(
	ctx context.Context, userID int, balance, withdrawal decimal.Decimal) error {
	return nil
}

func (r *GophermartRepo) AddWithdrawal(
	ctx context.Context, userID int, orderNum string, value decimal.Decimal) error {
	return nil
}

func (r *GophermartRepo) GetWithdrawalList(
	ctx context.Context, userID int, orderNum string) ([]entity.Withdraw, error) {
	return nil, nil
}
