package usecase

import (
	"context"
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vladislaoramos/gophermart/internal/entity"
	mocks "github.com/vladislaoramos/gophermart/internal/mocks/usecase"
	"github.com/vladislaoramos/gophermart/internal/repo"
	"github.com/vladislaoramos/gophermart/pkg/auth"
	"github.com/vladislaoramos/gophermart/pkg/logger"
	stdLog "log"
	"os"
	"testing"
	"time"
)

//type test struct {
//	name string
//	mock func()
//	arg  entity.User
//	res  interface{}
//	err  error
//}

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
	require.NoError(t, err)
	require.EqualValues(t, expectedBalance.Current, balance.Current)
	require.EqualValues(t, expectedBalance.Withdraw, balance.Withdraw)
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

func TestLoyalSystemUseCase_Withdraw(t *testing.T) {
	ls, repoMock := gophermart(t)
	ctx := context.Background()

	inputUserID := 1
	inputWithdrawal := entity.Withdrawal{Order: "123456", Sum: 100}
	expectedBalance := entity.Balance{Current: 500, Withdraw: 100}

	repoMock.On("GetBalance", mock.Anything, inputUserID).Return(expectedBalance, nil).Once()
	repoMock.On("UpdateBalance", mock.Anything, inputUserID, expectedBalance.Current-inputWithdrawal.Sum, expectedBalance.Withdraw+inputWithdrawal.Sum).Return(nil).Once()
	repoMock.On("AddWithdrawal", mock.Anything, inputUserID, inputWithdrawal.Order, inputWithdrawal.Sum).Return(nil).Once()

	err := ls.Withdraw(ctx, inputUserID, inputWithdrawal)
	require.NoError(t, err)

	repoMock.AssertExpectations(t)

	repoMock.AssertCalled(t, "GetBalance", ctx, inputUserID)
	repoMock.AssertCalled(t, "UpdateBalance", ctx, inputUserID, expectedBalance.Current-inputWithdrawal.Sum, expectedBalance.Withdraw+inputWithdrawal.Sum)
	repoMock.AssertCalled(t, "AddWithdrawal", ctx, inputUserID, inputWithdrawal.Order, inputWithdrawal.Sum)
}

func TestLoyalSystemUseCase_CheckUser(t *testing.T) {
	ls, mockRepo := gophermart(t)

	userAuth := entity.UserAuth{Login: "user", Password: "password"}
	hash, _ := auth.HashPassword(userAuth.Password)

	inputCtx := context.Background()
	expectedUser := entity.User{ID: 1, Login: userAuth.Login, PasswordHash: hash}
	mockRepo.On("GetUserWithLogin", inputCtx, userAuth.Login).Return(expectedUser, nil)

	user, err := ls.CheckUser(inputCtx, userAuth)

	require.Nil(t, err, "unexpected error")
	require.Equal(t, expectedUser, user, "unexpected user")

	mockRepo.AssertCalled(t, "GetUserWithLogin", inputCtx, userAuth.Login)
}

func TestLoyalSystemUseCase_CreateUser(t *testing.T) {
	ls, mockRepo := gophermart(t)

	ctx := context.Background()
	userAuth := entity.UserAuth{
		Login:    "testuser",
		Password: "testpassword",
	}

	mockRepo.On("GetUserWithLogin", mock.Anything, userAuth.Login).Return(entity.User{}, sql.ErrNoRows)

	user := entity.User{
		ID:    1,
		Login: userAuth.Login,
	}
	mockRepo.On("CreateUser", mock.Anything, userAuth.Login, mock.Anything).Return(user, nil)

	mockRepo.On("CreateUserBalance", mock.Anything, user.ID).Return(nil)

	result, err := ls.CreateUser(ctx, userAuth)
	require.NoError(t, err)
	require.EqualValues(t, user, result)

	mockRepo.AssertCalled(t, "GetUserWithLogin", mock.Anything, userAuth.Login)
	mockRepo.AssertCalled(t, "CreateUser", mock.Anything, userAuth.Login, mock.Anything)
	mockRepo.AssertCalled(t, "CreateUserBalance", mock.Anything, user.ID)
}

func TestUploadOrder(t *testing.T) {
	ls, repoMock := gophermart(t)

	inputUserID := 1
	inputOrderNum := "1234123412341234"

	repoMock.On("GetOrderByOrderNumber", mock.Anything, inputOrderNum).Return(entity.Order{}, repo.ErrNotFound)
	repoMock.On("CreateOrder", mock.Anything, inputUserID, inputOrderNum).Return(nil)
	repoMock.On("GetOrderByOrderNumber", mock.Anything, inputOrderNum).Return(entity.Order{UserID: inputUserID}, nil)

	_, err := ls.UploadOrder(context.Background(), inputUserID, inputOrderNum)
	require.NoError(t, err)

	repoMock.AssertCalled(t, "GetOrderByOrderNumber", mock.Anything, inputOrderNum)
	repoMock.AssertCalled(t, "CreateOrder", mock.Anything, inputUserID, inputOrderNum)
	repoMock.AssertCalled(t, "GetOrderByOrderNumber", mock.Anything, inputOrderNum)
}
