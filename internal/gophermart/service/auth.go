package service

import (
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/models"
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/storage"
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/utils"
)

type Auth interface {
	Register(login, password string) error
	Login(login, password string) (models.User, error)
}

var _ Auth = (*authService)(nil)

type authService struct {
	storage storage.Storage
}

func NewAuthService(s storage.Storage) *authService {
	return &authService{storage: s}
}

func (a *authService) Register(login, password string) error {
	securePass, err := utils.Generate(password)
	if err != nil {
		return err
	}

	return a.storage.CreateUser(login, securePass)
}

func (a *authService) Login(login, password string) (models.User, error) {
	user, err := a.storage.GetUserByName(login)

	if err != nil {
		return user, err
	}

	if ok := utils.Compare(password, user.Password); !ok {
		return user, utils.ErrUserPasswordMissMatch
	}

	return user, nil
}
