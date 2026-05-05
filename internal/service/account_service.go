package service

import (
	"fmt"

	"github.com/AndreLiberato/bitbank/internal/domain"
	"github.com/AndreLiberato/bitbank/internal/repository"
)

type AccountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) CreateAccount(number string) error {
	exists, err := s.repo.Exists(number)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("account %s already exists", number)
	}
	return s.repo.Save(domain.Account{Number: number, Balance: 0})
}

func (s *AccountService) GetBalance(number string) (float64, error) {
	account, err := s.repo.FindByNumber(number)
	if err != nil {
		return 0, err
	}
	if account == nil {
		return 0, fmt.Errorf("account %s not found", number)
	}
	return account.Balance, nil
}

func (s *AccountService) Credit(number string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	account, err := s.repo.FindByNumber(number)
	if err != nil {
		return err
	}
	if account == nil {
		return fmt.Errorf("account %s not found", number)
	}
	account.Balance += amount
	return s.repo.Update(*account)
}

func (s *AccountService) Debit(number string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	account, err := s.repo.FindByNumber(number)
	if err != nil {
		return err
	}
	if account == nil {
		return fmt.Errorf("account %s not found", number)
	}
	account.Balance -= amount
	return s.repo.Update(*account)
}

func (s *AccountService) Transfer(from, to string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	fromAccount, err := s.repo.FindByNumber(from)
	if err != nil {
		return err
	}
	if fromAccount == nil {
		return fmt.Errorf("source account %s not found", from)
	}
	toAccount, err := s.repo.FindByNumber(to)
	if err != nil {
		return err
	}
	if toAccount == nil {
		return fmt.Errorf("destination account %s not found", to)
	}
	fromAccount.Balance -= amount
	toAccount.Balance += amount
	if err := s.repo.Update(*fromAccount); err != nil {
		return err
	}
	return s.repo.Update(*toAccount)
}
