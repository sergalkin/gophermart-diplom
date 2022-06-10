package service

import (
	"errors"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/joeljunstrom/go-luhn"

	"github.com/sergalkin/gophermart-diplom/internal/gophermart/models"
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/storage"
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/utils"
)

type Orders interface {
	Create(login, order string) error
	Get(login string) ([]models.Order, error)
}

var _ Orders = (*ordersService)(nil)

type ordersService struct {
	storage storage.Storage
}

func NewOrdersService(s storage.Storage) *ordersService {
	return &ordersService{storage: s}
}

func (o *ordersService) Create(login, order string) error {
	if !luhn.Valid(order) {
		return utils.ErrLuhnValidation
	}

	if err := o.storage.CreateOrder(login, order); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if strings.Contains(pgErr.Message, "orders_order_idx") {
				return utils.ErrOrderAlreadyCreatedByOtherUser
			}
			if strings.Contains(pgErr.Message, "orders_order_user_idx") {
				return utils.ErrOrderUniqueViolation
			}
		}

		return err
	}

	return nil
}

func (o *ordersService) Get(login string) ([]models.Order, error) {
	orders, err := o.storage.GetOrders(login)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, utils.ErrOrdersNotFound
	}

	return orders, err
}
