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

func (s *AccountService) GetBalance(number string) (float64, error) {
	account, err := s.repo.FindByNumber(number)
	if err != nil {
		return 0, err
	}
	if account == nil {
		return 0, fmt.Errorf("conta %s não encontrada", number)
	}
	return account.Balance, nil
}

func (s *AccountService) CreateAccount(number string) error {
	exists, err := s.repo.Exists(number)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("conta %s já existe", number)
	}
	return s.repo.Save(domain.Account{Number: number, Balance: 0})
}
