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
		return fmt.Errorf("conta %s já existe", number)
	}
	return s.repo.Save(domain.Account{Number: number, Balance: 0})
}
