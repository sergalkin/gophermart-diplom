package storage

import (
	"context"
	"errors"
	"sync"

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
	createUser = `INSERT INTO public.users ("name", "password") VALUES ($1,$2)`
	getUser    = `SELECT * from public.users where "name" = $1`
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
