package service

import "github.com/AndreLiberato/bitbank/internal/repository"

type AccountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}
