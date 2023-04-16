package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/vladislaoramos/gophemart/internal/entity"
	"github.com/vladislaoramos/gophemart/internal/repo"
	"github.com/vladislaoramos/gophemart/internal/webapi"
	"github.com/vladislaoramos/gophemart/pkg/auth"
	"github.com/vladislaoramos/gophemart/pkg/logger"
	"time"
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
				err := ls.processOrder(orderNumber)
				if err != nil {
					log.Error(fmt.Errorf("UseCase - processOrder: %w", err).Error())
				}
			}
		}()
	}

	return ls
}

const (
	orderStatusProcessing = "PROCESSING"
	orderStatusProcessed  = "PROCESSED"
	orderStatusInvalid    = "INVALID"
)

func (ls *LoyalSystemUseCase) processOrder(orderNumber string) error {
	order, timeout, err := ls.webAPI.GetOrderInfo(orderNumber)

	switch err {
	case webapi.ErrTooManyRequests:
		time.Sleep(timeout)
		ls.orderCh <- orderNumber
		return fmt.Errorf("UseCase - processOrder - ls.webAPI.GetOrderInfo: %w", err)
	case webapi.ErrInternalServerError:
		err = ls.repo.UpdateOrderStatus(context.TODO(), orderNumber, orderStatusInvalid)
		if err != nil {
			return fmt.Errorf("UseCase - processOrder - ls.repo.UpdateOrderStatus: %w", err)
		}
		return err
	default:
		if err != nil {
			return fmt.Errorf("UseCase - processOrder: %w", err)
		}
	}

	addOrder, err := ls.repo.GetOrderByOrderNumber(context.TODO(), orderNumber)
	if err != nil {
		return fmt.Errorf("UseCase - ProcessOrder - ls.repo.GetOrderByOrderNumber: %w", err)
	}

	order.UserID = addOrder.UserID

	switch order.Status {
	case orderStatusInvalid:
		err = ls.repo.UpdateOrderStatus(context.Background(), orderNumber, order.Status)
		if err != nil {
			return fmt.Errorf("UseCase - ProcessOrder - ls.repo.UpdateOrderStatus: %w", err)
		}
		return nil
	case orderStatusProcessing:
		err = ls.repo.UpdateOrderStatus(context.Background(), orderNumber, order.Status)
		if err != nil {
			return fmt.Errorf("UseCase - ProcessOrder - ls.repo.UpdateOrderStatus: %w", err)
		}
		ls.orderCh <- orderNumber
	case orderStatusProcessed:
		err = ls.repo.UpdateOrderStatus(context.Background(), orderNumber, order.Status)
		if err != nil {
			return fmt.Errorf("UseCase - ProcessOrder - ls.repo.UpdateOrderStatus: %w", err)
		}
		err = ls.repo.UpdateOrderAccrual(context.Background(), orderNumber, order.Accrual)
		if err != nil {
			return fmt.Errorf("UseCase - ProcessOrder - ls.repo.UpdateOrderAccrual: %w", err)
		}
		balance, err := ls.repo.GetBalance(context.Background(), order.UserID)
		if err != nil {
			return fmt.Errorf("UseCase - ProcessOrder - ls.repo.GetCurrentBalance: %w", err)
		}
		err = ls.repo.UpdateBalance(context.Background(), order.UserID, balance.Current+order.Accrual, balance.Withdraw)
		if err != nil {
			return fmt.Errorf("UseCase - ProcessOrder - ls.repo.UpdateBalance: %w", err)
		}
	default:
		return fmt.Errorf("UseCase - processOrder: %w", err)
	}

	return nil
}

func (ls *LoyalSystemUseCase) CreateUser(
	ctx context.Context,
	userAuth entity.UserAuth,
) (entity.User, error) {
	pwdHash, err := auth.HashPassword(userAuth.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf("UseCase - CreateUser - auth.HashPassword: %w", err)
	}

	_, err = ls.repo.GetUserWithLogin(ctx, userAuth.Login)
	if err == nil {
		return entity.User{}, fmt.Errorf("UseCase - CreateUser: the user already exists: %w", ErrConflict)
	}

	user, err := ls.repo.CreateUser(ctx, userAuth.Login, pwdHash)
	if err != nil {
		return entity.User{}, fmt.Errorf("UseCase - CreateUser - ls.repo.CreateUser: %w", err)
	}

	err = ls.repo.CreateUserBalance(ctx, user.ID)
	if err != nil {
		return entity.User{}, fmt.Errorf("UseCase - CreateUser - ls.repo.CreateUserBalance: %w", err)
	}

	return user, nil
}

func (ls *LoyalSystemUseCase) CheckUser(ctx context.Context, userAuth entity.UserAuth) (entity.User, error) {
	pwdHash, err := auth.HashPassword(userAuth.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf("UseCase - CheckUser - auth.HashPassword: %w", err)
	}

	user, err := ls.repo.GetUserWithLogin(ctx, userAuth.Login)
	if err != nil {
		return entity.User{}, fmt.Errorf("UseCase - CheckUser - ls.repo.GetUserWithLogin: %w", ErrUnauthorized)
	}

	isValid := auth.ValidatePassword(userAuth.Password, pwdHash)
	if !isValid {
		return entity.User{}, fmt.Errorf("UseCase - CheckUser - auth.ValidatePassword: %w", ErrUnauthorized)
	}

	return user, nil
}

func (ls *LoyalSystemUseCase) PingRepo(ctx context.Context) error {
	return ls.repo.Ping(ctx)
}

func (ls *LoyalSystemUseCase) UploadOrder(
	ctx context.Context,
	userID int,
	orderNum string,
) (bool, error) {
	order, err := ls.repo.GetOrderByOrderNumber(ctx, orderNum)
	if err == nil {
		if order.UserID == userID {
			return true, nil
		}
		return false, fmt.Errorf("UseCase - UploadOrder - ls.repo.GetOrderByOrderNumber: %w", ErrConflict)
	}

	if errors.Is(err, repo.ErrNotFound) {
		err = ls.repo.CreateOrder(ctx, userID, orderNum)
		if err != nil {
			return false, fmt.Errorf("UseCase - UploadOrder - ls.repo.CreateOrder: %w", err)
		}
		ls.orderCh <- orderNum
	}

	if err != nil {
		return false, fmt.Errorf("UseCase - UploadOrder - ls.repo.GetOrderByOrderNumber: %w", err)
	}

	return false, nil
}

func (ls *LoyalSystemUseCase) GetOrderList(ctx context.Context, userID int) ([]entity.Order, error) {
	orders, err := ls.repo.GetOrderList(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("UseCase - GetOrderList - ls.repo.GetOrderList: %w", err)
	}
	return orders, nil
}

func (ls *LoyalSystemUseCase) GetBalance(ctx context.Context, userID int) (entity.Balance, error) {
	balance, err := ls.repo.GetBalance(ctx, userID)
	if err != nil {
		return entity.Balance{}, fmt.Errorf("UseCase - GetBalance - ls.repo.GetBalance: %w", err)
	}
	return balance, nil
}

func (ls *LoyalSystemUseCase) Withdraw(ctx context.Context, userID int, withdrawal entity.Withdrawal) error {
	balance, err := ls.repo.GetBalance(ctx, userID)
	if err != nil {
		return fmt.Errorf("UseCase - Withdraw - ls.repo.GetBalance: %w", err)
	}

	if withdrawal.Sum > balance.Current {
		return fmt.Errorf("UseCase - Withdraw: %w", ErrPaymentRequired)
	}

	err = ls.repo.UpdateBalance(ctx, userID, balance.Current-withdrawal.Sum, balance.Withdraw+withdrawal.Sum)
	if err != nil {
		return fmt.Errorf("UseCase - Withdraw - ls.repo.UpdateBalance: %w", err)
	}

	err = ls.repo.AddWithdrawal(ctx, userID, withdrawal.Order, withdrawal.Sum)
	if err != nil {
		return fmt.Errorf("UseCase - Withdraw - ls.repo.AddWithdrawal: %w", err)
	}

	return nil
}

func (ls *LoyalSystemUseCase) GetWithdrawList(ctx context.Context, userID int) ([]entity.Withdraw, error) {
	withdraw, err := ls.repo.GetWithdrawalList(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("UseCase - GetWithdrawList - ls.repo.GetWithdrawList: %w", err)
	}
	return withdraw, nil
}
