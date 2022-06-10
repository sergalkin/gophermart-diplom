package service

import (
	"github.com/joeljunstrom/go-luhn"
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/models"
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/storage"
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/utils"
)

type Balance interface {
	Get(login string) (models.Balance, error)
	Withdraw(sum float32, login, order string) error
	Withdrawals(login string) ([]models.Withdraw, error)
}

var _ Balance = (*balanceService)(nil)

type balanceService struct {
	storage storage.Storage
}

func NewBalanceService(s storage.Storage) *balanceService {
	return &balanceService{storage: s}
}

func (b *balanceService) Get(login string) (models.Balance, error) {
	return b.storage.GetBalance(login)
}

func (b *balanceService) Withdraw(sum float32, login, order string) error {
	if !luhn.Valid(order) {
		return utils.ErrLuhnValidation
	}

	if err := b.storage.Withdraw(sum, login, order); err != nil {
		return err
	}

	return nil
}

func (b *balanceService) Withdrawals(login string) ([]models.Withdraw, error) {
	w, err := b.storage.GetWithdrawals(login)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, utils.ErrEmptyRow
		}

		return nil, err
	}

	if len(w) == 0 {
		return nil, utils.ErrEmptyRow
	}

	return w, nil
}
