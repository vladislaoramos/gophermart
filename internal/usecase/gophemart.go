package usecase

import (
	"context"
	"github.com/vladislaoramos/gophemart/internal/entity"
	"github.com/vladislaoramos/gophemart/pkg/auth"
	"github.com/vladislaoramos/gophemart/pkg/logger"
)

type LoyalSystemUseCase struct {
	repo    LoyalSystemRepo
	webAPI  LoyalSystemAPI
	orderCh chan<- string
}

func NewLoyalSystem(
	r LoyalSystemRepo,
	w LoyalSystemAPI,
	workersNum int,
	log logger.LogInterface,
) *LoyalSystemUseCase {
	orders := make(chan string)
	ls := &LoyalSystemUseCase{
		repo:    r,
		webAPI:  w,
		orderCh: orders,
	}

	for i := 0; i < workersNum; i++ {
		go func() {
			for orderNumber := range orders {
				err := ls.ProcessOrder(orderNumber)
				if err != nil {
					log.Error(err.Error())
				}
			}
		}()
	}

	return ls
}

func (ls *LoyalSystemUseCase) CreateUser(
	ctx context.Context,
	userAuth entity.UserAuth,
) (entity.User, error) {
	pwdHash, err := auth.HashPassword(userAuth.Password)
	if err != nil {
		return entity.User{}, err
	}

	_, err = ls.repo.GetUserWithLogin(ctx, userAuth.Login)
	if err == nil {
		return entity.User{}, ErrConflict
	}

	user, err := ls.repo.CreateUser(ctx, userAuth.Login, pwdHash)
	if err != nil {
		return entity.User{}, err
	}

	err = ls.repo.CreateUserBalance(ctx, user.ID)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (ls *LoyalSystemUseCase) CheckUser(ctx context.Context, userAuth entity.UserAuth) (entity.User, error) {
	pwdHash, err := auth.HashPassword(userAuth.Password)
	if err != nil {
		return entity.User{}, err
	}

	user, err := ls.repo.GetUserWithLogin(ctx, userAuth.Login)
	if err != nil {
		return entity.User{}, ErrUnauthorized
	}

	isValid := auth.ValidatePassword(userAuth.Password, pwdHash)
	if !isValid {
		return entity.User{}, ErrUnauthorized
	}

	return user, nil
}

func (ls *LoyalSystemUseCase) ProcessOrder(orderNumber string) error {
	return nil
}

func (ls *LoyalSystemUseCase) PingRepo(ctx context.Context) error {
	return ls.repo.Ping(ctx)
}

func (ls *LoyalSystemUseCase) UploadOrder(ctx context.Context, userID int, orderNum string) (bool, error) {
	return false, nil
}

func (ls *LoyalSystemUseCase) GetOrderList(ctx context.Context, userID int) ([]entity.Order, error) {
	return nil, nil
}

func (ls *LoyalSystemUseCase) GetBalance(ctx context.Context, userID int) (entity.Balance, error) {
	return entity.Balance{}, nil
}

func (ls *LoyalSystemUseCase) Withdraw(ctx context.Context, userID int, withdrawal entity.Withdrawal) error {
	return nil
}

func (ls *LoyalSystemUseCase) GetWithdrawList(ctx context.Context, userID int) ([]entity.Withdraw, error) {
	return nil, nil
}
