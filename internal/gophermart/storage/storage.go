package storage

import "github.com/sergalkin/gophermart-diplom/internal/gophermart/models"

type Storage interface {
	CreateUser(login, password string) error
	GetUserByName(login string) (models.User, error)
	CreateOrder(login, order string) error
	GetOrders(login string) ([]models.Order, error)
	GetBalance(login string) (models.Balance, error)
	Withdraw(sum float32, login, order string) error
	UpdateBalance(login string, accrual float32) error
	GetWithdrawals(login string) ([]models.Withdraw, error)
	UpdateOrder(accrual float32, order, status string) error
}

func NewStorage() Storage {
	return NewDatabase()
}
