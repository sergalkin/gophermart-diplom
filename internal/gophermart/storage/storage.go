package storage

import "github.com/sergalkin/gophermart-diplom/internal/gophermart/models"

type Storage interface {
	CreateUser(login, password string) error
	GetUserByName(login string) (models.User, error)
	CreateOrder(login, order string) error
	GetOrders(login string) ([]models.Order, error)
}

func NewStorage() Storage {
	return NewDatabase()
}
