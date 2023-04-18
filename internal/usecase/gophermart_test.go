package usecase

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/vladislaoramos/gophermart/internal/entity"
	mocks "github.com/vladislaoramos/gophermart/internal/mocks/usecase"
	"github.com/vladislaoramos/gophermart/pkg/logger"
	stdLog "log"
	"os"
	"testing"
	"time"
)

type test struct {
	name string
	mock func()
	arg  entity.User
	res  interface{}
	err  error
}

func testLogger() *logger.Logger {
	f, err := os.OpenFile("/tmp/test_log_server", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	if err != nil {
		stdLog.Fatal("unable to open file for log")
	}

	return logger.New("debug", f)
}

func gophermart(t *testing.T) (*LoyalSystemUseCase, *mocks.LoyalSystemRepo) {
	log := testLogger()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockRepo := mocks.NewLoyalSystemRepo(t)
	mockWebAPI := mocks.NewLoyalSystemAPI(t)
	const workers = 5

	mockLS := NewLoyalSystem(
		mockRepo,
		mockWebAPI,
		workers,
		log,
	)

	return mockLS, mockRepo
}

func TestLoyalSystemUseCase_GetBalance(t *testing.T) {
	ls, mockRepo := gophermart(t)
	ctx := context.Background()

	userID := 1

	expectedBalance := entity.Balance{
		Current:  100.0,
		Withdraw: 50.0,
	}

	mockRepo.EXPECT().GetBalance(ctx, userID).Return(expectedBalance, nil)

	balance, err := ls.GetBalance(ctx, userID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if balance.Current != expectedBalance.Current || balance.Withdraw != expectedBalance.Withdraw {
		t.Fatalf("unexpected balance: %+v", balance)
	}
}

func TestLoyalSystemUseCase_GetWithdrawList(t *testing.T) {
	ls, mockRepo := gophermart(t)
	ctx := context.Background()
	userID := 1

	mockRepo.EXPECT().GetWithdrawalList(ctx, userID).Return(nil, nil)

	_, err := ls.GetWithdrawList(context.Background(), userID)
	require.NoError(t, err)
}

func TestLoyalSystemUseCase_GetOrderList(t *testing.T) {
	ls, mockRepo := gophermart(t)

	userID := 1
	expectedOrders := []entity.Order{
		{
			Number:     "123",
			Status:     orderStatusProcessed,
			Accrual:    10.0,
			UploadedAt: time.Now(),
			UserID:     userID,
		},
		{
			Number:     "456",
			Status:     orderStatusProcessed,
			Accrual:    5.0,
			UploadedAt: time.Now(),
			UserID:     userID,
		},
	}

	mockRepo.EXPECT().GetOrderList(context.Background(), userID).Return(expectedOrders, nil)

	orders, err := ls.GetOrderList(context.Background(), userID)
	require.NoError(t, err)
	require.Equal(t, expectedOrders, orders)
}
