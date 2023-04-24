package repo

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/vladislaoramos/gophermart/internal/entity"
	mocks "github.com/vladislaoramos/gophermart/internal/mocks/usecase"
	"testing"
)

//type fixture struct {
//	t          *testing.T
//	ctx        context.Context
//	repository *LoyalSystemRepo
//	db         *sqlx.DB
//}

//func tearUp(t *testing.T) *fixture {
//	dsn, ok := os.LookupEnv("DATABASE_URI")
//	require.True(t, ok)
//
//	db, err := postgres.New(dsn)
//	if err != nil {
//		log.Fatal(fmt.Errorf("app - Run - postgres.NewAPI: %w", err).Error())
//	}
//	defer db.Close()
//
//	return &fixture{
//		t:          t,
//		ctx:        context.Background(),
//		repository: NewRepository(db),
//	}
//}
//
//func (fx *fixture) tearDown() {
//	fx.repository.Close()
//}
//
//func (fx *fixture) loadFixtures(paths ...string) {
//
//}

//func (fx *fixture) updateOrderStatus(ctx context.Context, orderNumber, status string) error {
//	query, args, err := sq.Select("*").
//		From("public.orders").
//		Where(sq.Eq{"order_number": orderNumber}).
//		PlaceholderFormat(sq.Dollar).
//		ToSql()
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func TestUpdateOrderAccrual(t *testing.T) {
	t.Run("without error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := mocks.NewLoyalSystemRepo(t)

		orderNumber := "test_order"
		accrual := 100.0
		ctx := context.Background()

		mockRepo.EXPECT().UpdateOrderAccrual(
			ctx,
			orderNumber,
			accrual,
		).Return(nil)

		err := mockRepo.UpdateOrderAccrual(ctx, orderNumber, accrual)
		require.NoError(t, err)
	})

	t.Run("with error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := mocks.NewLoyalSystemRepo(t)

		orderNumber := "test_order"
		accrual := 100.0
		ctx := context.Background()

		mockRepo.EXPECT().UpdateOrderAccrual(
			ctx,
			orderNumber,
			accrual,
		).Return(fmt.Errorf("not found"))

		err := mockRepo.UpdateOrderAccrual(ctx, orderNumber, accrual)
		require.Error(t, err)
	})
}

func TestUpdateOrderStatus(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := mocks.NewLoyalSystemRepo(t)

		orderNumber := "test_order"
		status := "NEW"
		ctx := context.Background()

		mockRepo.EXPECT().UpdateOrderStatus(
			ctx,
			orderNumber,
			status,
		).Return(nil)

		err := mockRepo.UpdateOrderStatus(ctx, orderNumber, status)
		require.NoError(t, err)
	})

	t.Run("error too long status value", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := mocks.NewLoyalSystemRepo(t)

		orderNumber := "test_order"
		status := "NEWNEWNEWNEWNEWNEWNEWNEWNEWNEWNEWNEWNEW"
		ctx := context.Background()

		expectedErr := fmt.Errorf("too long status value")
		mockRepo.EXPECT().UpdateOrderStatus(
			ctx,
			orderNumber,
			status,
		).Return(expectedErr)

		err := mockRepo.UpdateOrderStatus(ctx, orderNumber, status)
		require.Error(t, err)
	})
}

func TestGetOrderByOrderNumber(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := mocks.NewLoyalSystemRepo(t)

		orderNumber := "1234567812345678"
		ctx := context.Background()

		mockRepo.EXPECT().GetOrderByOrderNumber(
			ctx,
			orderNumber,
		).Return(entity.Order{}, nil)

		_, err := mockRepo.GetOrderByOrderNumber(ctx, orderNumber)
		require.NoError(t, err)
	})

	t.Run("error order number", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := mocks.NewLoyalSystemRepo(t)

		orderNumber := "-"
		ctx := context.Background()

		expectedErr := fmt.Errorf("order number is invalid")
		mockRepo.EXPECT().GetOrderByOrderNumber(
			ctx,
			orderNumber,
		).Return(entity.Order{}, expectedErr)

		_, err := mockRepo.GetOrderByOrderNumber(ctx, orderNumber)
		require.Error(t, err)
	})

	t.Run("error order number", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := mocks.NewLoyalSystemRepo(t)

		orderNumber := "1234567812345679"
		ctx := context.Background()

		expectedErr := fmt.Errorf("order number is not found")
		mockRepo.EXPECT().GetOrderByOrderNumber(
			ctx,
			orderNumber,
		).Return(entity.Order{}, expectedErr)

		_, err := mockRepo.GetOrderByOrderNumber(ctx, orderNumber)
		require.Error(t, err)
	})
}

func TestCreateUserBalance(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := mocks.NewLoyalSystemRepo(t)

		userID := 1
		ctx := context.Background()

		mockRepo.EXPECT().CreateUserBalance(
			ctx,
			userID,
		).Return(nil)

		err := mockRepo.CreateUserBalance(ctx, userID)
		require.NoError(t, err)
	})

	t.Run("error invalid user ID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := mocks.NewLoyalSystemRepo(t)

		userID := 0
		ctx := context.Background()

		expectedErr := fmt.Errorf("invalid user id")
		mockRepo.EXPECT().CreateUserBalance(
			ctx,
			userID,
		).Return(expectedErr)

		err := mockRepo.CreateUserBalance(ctx, userID)
		require.Error(t, err)
	})
}

func TestCreateOrder(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := mocks.NewLoyalSystemRepo(t)

		userID := 1
		orderNum := "12345"
		ctx := context.Background()

		mockRepo.EXPECT().CreateOrder(
			ctx,
			userID,
			orderNum,
		).Return(nil)

		err := mockRepo.CreateOrder(ctx, userID, orderNum)
		require.NoError(t, err)
	})
}

func TestCreateUser(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := mocks.NewLoyalSystemRepo(t)

		login := "user-login"
		pwdHash := "pwd-secret-hash"
		ctx := context.Background()

		mockRepo.EXPECT().CreateUser(
			ctx,
			login,
			pwdHash,
		).Return(entity.User{Login: login, PasswordHash: pwdHash}, nil)

		user, err := mockRepo.CreateUser(ctx, login, pwdHash)
		require.NoError(t, err)
		require.Equal(t, login, user.Login)
		require.Equal(t, pwdHash, user.PasswordHash)
	})
}

func TestGetUserWithLogin(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := mocks.NewLoyalSystemRepo(t)

		login := "user-login"
		ctx := context.Background()

		mockRepo.EXPECT().GetUserWithLogin(
			ctx,
			login,
		).Return(entity.User{Login: login}, nil)

		user, err := mockRepo.GetUserWithLogin(ctx, login)
		require.NoError(t, err)
		require.Equal(t, login, user.Login)
	})

	t.Run("error user is not found", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := mocks.NewLoyalSystemRepo(t)

		login := "not-found-user-login"
		ctx := context.Background()

		expectedErr := fmt.Errorf("not found")
		mockRepo.EXPECT().GetUserWithLogin(
			ctx,
			login,
		).Return(entity.User{}, expectedErr)

		_, err := mockRepo.GetUserWithLogin(ctx, login)
		require.Error(t, err)
	})
}

func TestGetOrderList(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := mocks.NewLoyalSystemRepo(t)

		userID := 1
		ctx := context.Background()

		mockRepo.EXPECT().GetOrderList(
			ctx,
			userID,
		).Return([]entity.Order{}, nil)

		orders, err := mockRepo.GetOrderList(ctx, userID)
		require.NoError(t, err)
		require.NotNil(t, orders)
	})

	t.Run("invalid user id", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := mocks.NewLoyalSystemRepo(t)

		userID := 0
		ctx := context.Background()

		expectedErr := fmt.Errorf("invalid user id")
		mockRepo.EXPECT().GetOrderList(
			ctx,
			userID,
		).Return(nil, expectedErr)

		orders, err := mockRepo.GetOrderList(ctx, userID)
		require.Error(t, err)
		require.Nil(t, orders)
	})
}

func TestGetBalance(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := mocks.NewLoyalSystemRepo(t)

		userID := 1
		ctx := context.Background()

		mockRepo.EXPECT().GetBalance(
			ctx,
			userID,
		).Return(entity.Balance{Current: 50.0}, nil)

		balance, err := mockRepo.GetBalance(ctx, userID)
		require.NoError(t, err)
		require.Equal(t, 50.0, balance.Current)
	})

	t.Run("invalid user id", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := mocks.NewLoyalSystemRepo(t)

		userID := 0
		ctx := context.Background()

		expectedErr := fmt.Errorf("invalid user id")
		mockRepo.EXPECT().GetBalance(
			ctx,
			userID,
		).Return(entity.Balance{}, expectedErr)

		_, err := mockRepo.GetBalance(ctx, userID)
		require.Error(t, err)
	})
}

func TestGetWithdrawalList(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := mocks.NewLoyalSystemRepo(t)

		userID := 1
		ctx := context.Background()

		mockRepo.EXPECT().GetWithdrawalList(
			ctx,
			userID,
		).Return([]entity.Withdraw{}, nil)

		withdraws, err := mockRepo.GetWithdrawalList(ctx, userID)
		require.NoError(t, err)
		require.NotNil(t, withdraws)
	})

	t.Run("invalid user id", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := mocks.NewLoyalSystemRepo(t)

		userID := 0
		ctx := context.Background()

		expectedErr := fmt.Errorf("invalid user id")
		mockRepo.EXPECT().GetWithdrawalList(
			ctx,
			userID,
		).Return(nil, expectedErr)

		withdraws, err := mockRepo.GetWithdrawalList(ctx, userID)
		require.Error(t, err)
		require.Nil(t, withdraws)
	})
}

func TestUpdateBalance(t *testing.T) {
	t.Run("without error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := mocks.NewLoyalSystemRepo(t)

		balance := 100.0
		withdrawal := 50.0
		userID := 1
		ctx := context.Background()

		mockRepo.EXPECT().UpdateBalance(
			ctx,
			userID,
			balance,
			withdrawal,
		).Return(nil)

		err := mockRepo.UpdateBalance(ctx, userID, balance, withdrawal)
		require.NoError(t, err)
	})

	t.Run("error user not found", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := mocks.NewLoyalSystemRepo(t)

		balance := 100.0
		withdrawal := 50.0
		userID := 0
		ctx := context.Background()

		expectedErr := fmt.Errorf("user is not found")
		mockRepo.EXPECT().UpdateBalance(
			ctx,
			userID,
			balance,
			withdrawal,
		).Return(expectedErr)

		err := mockRepo.UpdateBalance(ctx, userID, balance, withdrawal)
		require.Error(t, err)
	})
}

func TestAddWithdrawal(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		mockRepo := mocks.NewLoyalSystemRepo(t)

		userID := 1
		orderNum := "12345"
		value := 100.0
		ctx := context.Background()

		mockRepo.On("AddWithdrawal", ctx, userID, orderNum, value).Return(nil)

		err := mockRepo.AddWithdrawal(ctx, userID, orderNum, value)

		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error negative user ID", func(t *testing.T) {
		mockRepo := mocks.NewLoyalSystemRepo(t)

		userID := 0
		orderNum := "12345"
		value := 100.0
		ctx := context.Background()

		expectedError := fmt.Errorf("negative user ID")
		mockRepo.On("AddWithdrawal", ctx, userID, orderNum, value).Return(expectedError)

		err := mockRepo.AddWithdrawal(ctx, userID, orderNum, value)

		require.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
