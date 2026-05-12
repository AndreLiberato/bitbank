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

func (s *AccountService) Transfer(from, to string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("valor deve ser maior que zero")
	}
	fromAccount, err := s.repo.FindByNumber(from)
	if err != nil {
		return err
	}
	if fromAccount == nil {
		return fmt.Errorf("conta origem %s não encontrada", from)
	}
	toAccount, err := s.repo.FindByNumber(to)
	if err != nil {
		return err
	}
	if toAccount == nil {
		return fmt.Errorf("conta destino %s não encontrada", to)
	}
	if fromAccount.Balance < amount {
		return fmt.Errorf("saldo insuficiente na conta origem %s", from)
	}
	fromAccount.Balance -= amount
	toAccount.Balance += amount
	if toAccount.IsBonus() {
		toAccount.Points += pointsForReceivedTransfer(amount)
	}
	if err := s.repo.Update(*fromAccount); err != nil {
		return err
	}
	return s.repo.Update(*toAccount)
}

func (s *AccountService) Debit(number string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("valor deve ser maior que zero")
	}
	account, err := s.repo.FindByNumber(number)
	if err != nil {
		return err
	}
	if account == nil {
		return fmt.Errorf("conta %s não encontrada", number)
	}
	if account.Balance < amount {
		return fmt.Errorf("saldo insuficiente na conta %s", number)
	}
	account.Balance -= amount
	return s.repo.Update(*account)
}

func (s *AccountService) Credit(number string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("valor deve ser maior que zero")
	}
	account, err := s.repo.FindByNumber(number)
	if err != nil {
		return err
	}
	if account == nil {
		return fmt.Errorf("conta %s não encontrada", number)
	}
	account.Balance += amount
	if account.IsBonus() {
		account.Points += pointsForDeposit(amount)
	}
	return s.repo.Update(*account)
}

func (s *AccountService) GetBalance(number string) (float64, error) {
	account, err := s.GetAccount(number)
	if err != nil {
		return 0, err
	}
	return account.Balance, nil
}

func (s *AccountService) GetAccount(number string) (*domain.Account, error) {
	account, err := s.repo.FindByNumber(number)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, fmt.Errorf("conta %s não encontrada", number)
	}
	return account, nil
}

func (s *AccountService) CreateAccount(number string) error {
	return s.createAccount(domain.Account{
		Number:  number,
		Balance: 0,
		Type:    domain.AccountTypeSimple,
		Points:  0,
	})
}

func (s *AccountService) CreateBonusAccount(number string) error {
	return s.createAccount(domain.Account{
		Number:  number,
		Balance: 0,
		Type:    domain.AccountTypeBonus,
		Points:  10,
	})
}

func (s *AccountService) CreateSavingsAccount(number string) error {
	return s.createAccount(domain.Account{
		Number:  number,
		Balance: 0,
		Type:    domain.AccountTypeSavings,
		Points:  0,
	})
}

func (s *AccountService) RenderJuros(taxa float64) error {
	if taxa <= 0 {
		return fmt.Errorf("taxa deve ser maior que zero")
	}
	accounts, err := s.repo.FindAll()
	if err != nil {
		return err
	}
	for _, account := range accounts {
		if !account.IsSavings() {
			continue
		}
		account.Balance += account.Balance * (taxa / 100)
		if err := s.repo.Update(account); err != nil {
			return err
		}
	}
	return nil
}

func (s *AccountService) createAccount(account domain.Account) error {
	exists, err := s.repo.Exists(account.Number)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("conta %s já existe", account.Number)
	}
	return s.repo.Save(account)
}

func pointsForDeposit(amount float64) int {
	return int(amount / 100)
}

func pointsForReceivedTransfer(amount float64) int {
	return int(amount / 200)
}
