package repository

import (
	"fmt"

	"github.com/AndreLiberato/bitbank/internal/domain"
)

// memoryRepository é uma implementação em memória de AccountRepository.
// É útil para testes unitários da camada de serviço, permitindo exercitar
// as operações do banco sem depender do SQLite.
type memoryRepository struct {
	accounts map[string]domain.Account
}

// NewMemoryRepository cria um repositório em memória vazio.
func NewMemoryRepository() AccountRepository {
	return &memoryRepository{accounts: make(map[string]domain.Account)}
}

func (r *memoryRepository) FindAll() ([]domain.Account, error) {
	accounts := make([]domain.Account, 0, len(r.accounts))
	for _, a := range r.accounts {
		accounts = append(accounts, a)
	}
	return accounts, nil
}

func (r *memoryRepository) FindByNumber(number string) (*domain.Account, error) {
	a, ok := r.accounts[number]
	if !ok {
		return nil, nil
	}
	cp := a
	return &cp, nil
}

func (r *memoryRepository) Save(account domain.Account) error {
	if _, ok := r.accounts[account.Number]; ok {
		return fmt.Errorf("conta %s já existe", account.Number)
	}
	r.accounts[account.Number] = account
	return nil
}

func (r *memoryRepository) Update(account domain.Account) error {
	if _, ok := r.accounts[account.Number]; !ok {
		return fmt.Errorf("account %s not found", account.Number)
	}
	r.accounts[account.Number] = account
	return nil
}

func (r *memoryRepository) Exists(number string) (bool, error) {
	_, ok := r.accounts[number]
	return ok, nil
}
