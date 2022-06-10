package storage

import (
	"context"
	"errors"
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
	createUser  = `INSERT INTO public.users ("name", "password") VALUES ($1,$2)`
	getUser     = `SELECT * from public.users where "name" = $1`
	createOrder = `INSERT INTO public.orders ("order","name","uploaded_at") VALUES ($1,$2,$3)`
	getOrders   = `SELECT "order", "status", "accrual", "uploaded_at" from public.orders where "name" = $1 ORDER BY "uploaded_at" DESC`
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
