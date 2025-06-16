package service

import (
	"errors"
	"github.com/google/uuid"
	"log"
	"payments/infrastructure/storage"
	model "payments/models"
)

type AccountServicing interface {
	Add(string, string, float64) error
	Get(id string) (*model.Account, error)
	All() ([]*model.Account, error)
	Update(string, float64) error
}

type AccountService struct {
	accounts storage.AccountStorer
	lg       *log.Logger
}

func NewAccountService(accounts storage.AccountStorer, lg *log.Logger) *AccountService {
	return &AccountService{accounts: accounts, lg: lg}
}

func (s *AccountService) Add(userID string, fullname string, balance float64) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	if balance < 0 {
		return errors.New("balance must be greater than zero")
	}
	if fullname == "" {
		return errors.New("full name must not be empty")
	}

	return s.accounts.Add(model.NewAccount(uid, fullname, balance))
}

func (s *AccountService) Get(userID string) (*model.Account, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return s.accounts.Get(uid)
}

func (s *AccountService) All() ([]*model.Account, error) {
	return s.accounts.All()
}

func (s *AccountService) Update(userID string, amount float64) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	return s.accounts.Update(uid, amount)
}
