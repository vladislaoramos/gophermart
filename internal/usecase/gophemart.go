package usecase

import (
	"context"
	"github.com/vladislaoramos/gophemart/internal/entity"
)

type GophermartUseCase struct {
	repo    GophermartRepo
	webAPI  GophermartWebAPI
	orderCh chan<- string
}

func New(r GophermartRepo, w GophermartWebAPI, workersCount int) *GophermartUseCase {
	return nil
}

func (uc *GophermartUseCase) ProcessOrder(orderNumber string) error {
	return nil
}

func (uc *GophermartUseCase) PingRepo(ctx context.Context) error {
	return nil
}

func (uc *GophermartUseCase) CreateNewUser(ctx context.Context, userAuth entity.UserAuth) (entity.User, error) {
	return entity.User{}, nil
}

func (uc *GophermartUseCase) CheckUser(ctx context.Context, userAuth entity.UserAuth) (entity.User, error) {
	return entity.User{}, nil
}

func (uc *GophermartUseCase) UploadOrder(ctx context.Context, userID int, orderNum string) (bool, error) {
	return false, nil
}

func (uc *GophermartUseCase) GetOrderList(ctx context.Context, userID int) ([]entity.Order, error) {
	return nil, nil
}

func (uc *GophermartUseCase) GetCurrentBalance(ctx context.Context, userID int) (entity.Balance, error) {
	return entity.Balance{}, nil
}

func (uc *GophermartUseCase) Withdraw(ctx context.Context, userID int, withdrawal entity.Withdrawal) error {
	return nil
}

func (uc *GophermartUseCase) GetWithdrawList(ctx context.Context, userID int) ([]entity.Withdraw, error) {
	return nil, nil
}
