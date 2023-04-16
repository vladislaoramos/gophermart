package repo

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/vladislaoramos/gophemart/internal/entity"
	"github.com/vladislaoramos/gophemart/pkg/postgres"
)

type LoyalSystemRepo struct {
	*postgres.DB
}

func NewRepository(db *postgres.DB) *LoyalSystemRepo {
	return &LoyalSystemRepo{db}
}

func (r *LoyalSystemRepo) Ping(_ context.Context) error {
	return nil
}

func (r *LoyalSystemRepo) UpdateOrderAccrual(
	ctx context.Context, orderNumber string, accrual float64) error {
	query, args, err := r.Builder.
		Update("public.orders").
		Set("accrual", accrual).
		Where(sq.Eq{"order_number": orderNumber}).
		ToSql()
	if err != nil {
		return fmt.Errorf("LoyalSystemRepo - UpdateOrderAccrual - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("LoyalSystemRepo - UpdateOrderAccrual - r.Pool.Exec: %w", err)
	}
	return nil
}

func (r *LoyalSystemRepo) UpdateOrderStatus(
	ctx context.Context, orderNumber, status string) error {
	query, args, err := r.Builder.
		Update("public.orders").
		Set("status", status).
		Where(sq.Eq{"order_number": orderNumber}).
		ToSql()
	if err != nil {
		return fmt.Errorf("LoyalSystemRepo - UpdateOrderStatus - r.builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("LoyalSystemRepo - UpdateOrderStatus - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *LoyalSystemRepo) GetOrderByOrderNumber(
	ctx context.Context, orderNumber string) (entity.Order, error) {
	query, args, err := r.Builder.
		Select("order_number", "status", "accrual", "uploaded_at", "user_id").
		From("public.orders").
		Where(sq.Eq{"order_number": orderNumber}).
		OrderBy("uploaded_at").
		ToSql()
	if err != nil {
		return entity.Order{}, fmt.Errorf("LoyalSystemRepo - GetOrderByOrderNumber - r.Builder: %w", err)
	}

	dst := make([]entity.Order, 0)
	if err = pgxscan.Select(ctx, r.Pool, &dst, query, args...); err != nil {
		return entity.Order{}, fmt.Errorf("LoyalSystemRepo - GetOrderByOrderNumber - pgxsan.Select: %w", err)
	}

	if len(dst) == 0 {
		return entity.Order{}, ErrNotFound
	}

	return dst[0], nil
}

func (r *LoyalSystemRepo) CreateUserBalance(
	ctx context.Context, userID int) error {
	query, args, err := r.Builder.
		Insert("public.balance").
		Columns("balance", "withdrawal", "user_id").
		Values(0, 0, userID).
		ToSql()
	if err != nil {
		return fmt.Errorf("LoyalSystemRepo - CreateUserBalance - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("LoyalSystemRepo - CreateUserBalance - r.Pool.Exec: %w", err)
	}
	return nil
}

func (r *LoyalSystemRepo) CreateOrder(
	ctx context.Context, userID int, orderNumber string) error {
	query, args, err := r.Builder.
		Insert("public.orders").
		Columns("order_number", "user_id", "accrual").
		Values(orderNumber, userID, 0).
		ToSql()
	if err != nil {
		return fmt.Errorf("LoyalSystemRepo - CreateOrder - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("LoyalSystemRepo - CreateOrder - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *LoyalSystemRepo) CreateUser(
	ctx context.Context, login, pwdHash string) (entity.User, error) {
	user := entity.User{
		Login:        login,
		PasswordHash: pwdHash,
	}

	query, args, err := r.Builder.
		Insert("public.users").
		Columns("login", "password_hash").
		Values(login, pwdHash).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return entity.User{}, fmt.Errorf("LoyalSystemRepo - CreateUser - r.Builder: %w", err)
	}

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return entity.User{}, fmt.Errorf("LoyalSystemRepo - CreateUser - r.Pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, query, args...).Scan(&user.ID)
	if err != nil {
		return entity.User{}, fmt.Errorf("LoyalSystemRepo - CreateUser - tx.QueryRow: %w", err)
	}

	tx.Commit(ctx)

	return user, nil
}

func (r *LoyalSystemRepo) GetUserWithLogin(
	ctx context.Context, login string) (entity.User, error) {
	query, args, err := r.Builder.
		Select("id", "login", "password_hash").
		From("public.users").
		Where(sq.Eq{"login": login}).
		ToSql()

	if err != nil {
		return entity.User{}, fmt.Errorf("LoyalSystemRepo - GetUserWithLogin - r.Builder: %w", err)
	}

	dst := make([]entity.User, 0)
	if err = pgxscan.Select(ctx, r.Pool, &dst, query, args...); err != nil {
		return entity.User{}, fmt.Errorf("LoyalSystemRepo - GetUserWithLogin - pgxscan.Select: %w", err)
	}

	if len(dst) == 0 {
		return entity.User{}, ErrNotFound
	}

	return dst[0], nil
}

func (r *LoyalSystemRepo) GetOrderList(
	ctx context.Context, userID int) ([]entity.Order, error) {
	query, args, err := r.Builder.
		Select("order_number", "status", "accrual", "uploaded_at").
		From("public.orders").
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("LoyalSystemRepo - GetOrderList - r.Builder: %w", err)
	}

	dst := make([]entity.Order, 0)
	if err = pgxscan.Select(ctx, r.Pool, &dst, query, args...); err != nil {
		return nil, fmt.Errorf("LoyalSystemRepo - GetOrderList - pgxscan.Select: %w", err)
	}

	return dst, nil
}

func (r *LoyalSystemRepo) GetBalance(
	ctx context.Context, userID int) (entity.Balance, error) {
	query, args, err := r.Builder.
		Select("balance", "withdrawal").
		From("public.balance").
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return entity.Balance{}, fmt.Errorf("LoyalSystemRepo - GetBalance - r.Builder: %w", err)
	}

	dst := make([]entity.Balance, 0)
	if err = pgxscan.Select(ctx, r.Pool, &dst, query, args...); err != nil {
		return entity.Balance{}, fmt.Errorf("LoyalSystemRepo - GetBalance - pgxsan.Select: %w", err)
	}

	if len(dst) == 0 {
		return entity.Balance{}, ErrNotFound
	}

	return dst[0], nil
}

func (r *LoyalSystemRepo) UpdateBalance(
	ctx context.Context, userID int, balance, withdrawal float64) error {
	query, args, err := r.Builder.
		Update("public.balance").
		Set("balance", balance).
		Set("withdrawal", withdrawal).
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("LoyalSystemRepo - UpdateBalance - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("LoyalSystemRepo - UpdateBalance - r.Pool.Exec: %w", err)
	}
	return nil
}

func (r *LoyalSystemRepo) AddWithdrawal(
	ctx context.Context, userID int, orderNum string, value float64) error {
	query, args, err := r.Builder.
		Insert("public.withdrawal").
		Columns("order_number", "sum_number", "user_id").
		Values(orderNum, value, userID).
		ToSql()
	if err != nil {
		return fmt.Errorf("LoyalSystemRepo - AddWithdrawal - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("LoyalSystemRepo - AddWithdrawal - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *LoyalSystemRepo) GetWithdrawalList(
	ctx context.Context, userID int) ([]entity.Withdraw, error) {
	query, args, err := r.Builder.
		Select("order_number", "sum_number", "updated_at").
		From("public.withdrawal").
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("LoyalSystemRepo - GetWithdrawalList - r.Builder: %w", err)
	}

	dst := make([]entity.Withdraw, 0)
	if err = pgxscan.Select(ctx, r.Pool, &dst, query, args...); err != nil {
		return nil, fmt.Errorf("LoyalSystemRepo - GetWithdrawalList - pgxscan.Select: %w", err)
	}

	if len(dst) == 0 {
		return nil, nil
	}

	return dst, nil
}
