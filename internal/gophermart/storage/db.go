package storage

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/sergalkin/gophermart-diplom/internal/gophermart/config"
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/migrations"
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/models"
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/utils"
)

var (
	db   database
	once sync.Once
	_    Storage = (*database)(nil)
)

const (
	createUser        = `INSERT INTO public.users ("name", "password") VALUES ($1,$2)`
	getUser           = `SELECT * from public.users where "name" = $1`
	createOrder       = `INSERT INTO public.orders ("order","name","uploaded_at") VALUES ($1,$2,$3)`
	getOrders         = `SELECT "order", "status", "accrual", "uploaded_at" from public.orders where "name" = $1 ORDER BY "uploaded_at" DESC`
	getBalance        = `SELECT "balance", "withdraw" FROM public.balance where "name" = $1`
	makeWithdrawal    = `INSERT INTO public.withdrawals ("name", "order", "processed_at", "withdraw")  VALUES ($1,$2,$3,$4)`
	updateBalance     = `UPDATE public.balance SET balance=$1, withdraw=$2 WHERE "name" = $3`
	createUserBalance = `INSERT INTO public.balance ("name") VALUES ($1)`
	getWithdrawals    = `SELECT "order", "withdraw", "processed_at" FROM public.withdrawals WHERE "name" = $1 ORDER BY "processed_at" DESC`
	updateOrder       = `UPDATE public.orders SET status=$1, accrual=$2 WHERE "order" = $3`
)

type database struct {
	conn *pgxpool.Pool
}

func NewDatabase() *database {
	once.Do(func() {
		conn, err := pgxpool.Connect(context.Background(), config.DatabaseDSN())
		if err != nil {
			config.Logger().Fatal("Error in creating pgxpool connection")
		}

		db.conn = conn

		if _, err = migrations.Up(); err != nil {
			config.Logger().Fatal("Error in running migrations")
		}
	})

	return &db
}

func (d *database) CreateUser(l, p string) error {
	_, err := d.conn.Exec(context.Background(), createUser, l, p)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return utils.ErrUserAlreadyExists
		}

		return err
	}

	if _, err = d.conn.Exec(context.Background(), createUserBalance, l); err != nil {
		return err
	}

	return nil
}

func (d *database) GetUserByName(l string) (models.User, error) {
	user := models.User{}

	err := d.conn.QueryRow(context.Background(), getUser, l).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {

		if err.Error() == "no rows in result set" {
			return user, utils.ErrUserNotFound
		}

		return user, err
	}

	return user, nil
}

func (d *database) CreateOrder(l, order string) error {
	_, err := d.conn.Exec(context.Background(), createOrder, order, l, time.Now())

	return err
}

func (d *database) GetOrders(login string) ([]models.Order, error) {
	orders := make([]models.Order, 0)
	rows, err := d.conn.Query(context.Background(), getOrders, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		order := models.Order{}

		err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.Upload)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (d *database) GetBalance(l string) (models.Balance, error) {
	balance := models.Balance{}

	err := d.conn.QueryRow(context.Background(), getBalance, l).Scan(&balance.Balance, &balance.Withdraws)
	if err != nil {
		return balance, err
	}

	return balance, nil
}

func (d *database) Withdraw(sum float32, l, order string) error {
	sum = sum * -1
	err := d.UpdateBalance(l, sum)
	if err != nil {
		return err
	}

	_, err = d.conn.Exec(context.Background(), makeWithdrawal, l, order, time.Now(), float32(math.Abs(float64(sum))))
	if err != nil {
		return err
	}

	return nil
}

func (d *database) UpdateBalance(l string, accrual float32) error {
	balance, err := d.GetBalance(l)
	if err != nil {
		return err
	}

	balance.Balance += accrual
	if balance.Balance < 0 {
		return utils.ErrNegativeBalance
	}

	if accrual < 0 {
		balance.Withdraws += float32(math.Abs(float64(accrual)))
	}

	_, err = d.conn.Exec(context.Background(), updateBalance, balance.Balance, balance.Withdraws, l)
	if err != nil {
		return err
	}

	return nil
}

func (d *database) GetWithdrawals(l string) ([]models.Withdraw, error) {
	withdrawals := make([]models.Withdraw, 0)
	rows, err := d.conn.Query(context.Background(), getWithdrawals, l)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		withdraw := models.Withdraw{}

		err = rows.Scan(&withdraw.Number, &withdraw.Withdraw, &withdraw.Processed)
		if err != nil {
			return nil, err
		}

		if withdraw.Withdraw > 0 {
			withdrawals = append(withdrawals, withdraw)
		}
	}

	return withdrawals, nil
}

func (d *database) UpdateOrder(accrual float32, order, status string) error {
	if _, err := d.conn.Exec(context.Background(), updateOrder, status, accrual, order); err != nil {
		return err
	}

	return nil
}
