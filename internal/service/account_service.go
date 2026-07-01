package service

import (
	"fmt"

	"github.com/AndreLiberato/bitbank/internal/domain"
	"github.com/AndreLiberato/bitbank/internal/repository"
)

const negativeBalanceLimit = -1000.0

// AccountService é a fachada (camada de negócios) do banco. Concentra todas as
// regras de negócio e pode ser invocada diretamente — sem passar pela API REST.
// Depende apenas da abstração repository.AccountRepository.
type AccountService struct {
	repo repository.AccountRepository
}

// NewAccountService cria a fachada a partir de um repositório de contas.
func NewAccountService(repo repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

// OpenAccount cadastra uma conta de acordo com o tipo informado. É um ponto de
// entrada único usado pela API REST; delega para os métodos específicos.
func (s *AccountService) OpenAccount(number, accountType string, initialBalance float64) error {
	switch accountType {
	case domain.AccountTypeBonus:
		return s.CreateBonusAccount(number)
	case domain.AccountTypeSavings:
		return s.CreateSavingsAccount(number, initialBalance)
	case domain.AccountTypeSimple, "":
		return s.CreateAccount(number, initialBalance)
	default:
		return fmt.Errorf("%w: %s", ErrInvalidAccountType, accountType)
	}
}

// CreateAccount cadastra uma conta simples.
func (s *AccountService) CreateAccount(number string, initialBalance float64) error {
	if initialBalance < 0 {
		return ErrNegativeInitialBalance
	}
	return s.createAccount(domain.Account{
		Number:  number,
		Balance: initialBalance,
		Type:    domain.AccountTypeSimple,
		Points:  0,
	})
}

// CreateBonusAccount cadastra uma conta bônus (inicia com 10 pontos).
func (s *AccountService) CreateBonusAccount(number string) error {
	return s.createAccount(domain.Account{
		Number:  number,
		Balance: 0,
		Type:    domain.AccountTypeBonus,
		Points:  10,
	})
}

// CreateSavingsAccount cadastra uma conta poupança.
func (s *AccountService) CreateSavingsAccount(number string, initialBalance float64) error {
	if initialBalance < 0 {
		return ErrNegativeInitialBalance
	}
	return s.createAccount(domain.Account{
		Number:  number,
		Balance: initialBalance,
		Type:    domain.AccountTypeSavings,
		Points:  0,
	})
}

// GetAccount consulta os dados completos de uma conta.
func (s *AccountService) GetAccount(number string) (*domain.Account, error) {
	account, err := s.repo.FindByNumber(number)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, fmt.Errorf("%w: %s", ErrAccountNotFound, number)
	}
	return account, nil
}

// GetBalance consulta o saldo de uma conta.
func (s *AccountService) GetBalance(number string) (float64, error) {
	account, err := s.GetAccount(number)
	if err != nil {
		return 0, err
	}
	return account.Balance, nil
}

// Credit credita um valor na conta. Contas bônus acumulam pontos.
func (s *AccountService) Credit(number string, amount float64) error {
	if err := validateAmount(amount); err != nil {
		return err
	}
	account, err := s.GetAccount(number)
	if err != nil {
		return err
	}
	account.Balance += amount
	if account.IsBonus() {
		account.Points += pointsForDeposit(amount)
	}
	return s.repo.Update(*account)
}

// Debit debita um valor da conta, respeitando o limite de saldo do tipo.
func (s *AccountService) Debit(number string, amount float64) error {
	if err := validateAmount(amount); err != nil {
		return err
	}
	account, err := s.GetAccount(number)
	if err != nil {
		return err
	}
	newBalance := account.Balance - amount
	if !allowsBalance(*account, newBalance) {
		return fmt.Errorf("%w na conta %s", ErrInsufficientBalance, number)
	}
	account.Balance = newBalance
	return s.repo.Update(*account)
}

// Transfer transfere um valor entre duas contas. A conta destino, se bônus,
// acumula pontos pela transferência recebida.
func (s *AccountService) Transfer(from, to string, amount float64) error {
	if err := validateAmount(amount); err != nil {
		return err
	}
	fromAccount, err := s.GetAccount(from)
	if err != nil {
		return err
	}
	toAccount, err := s.GetAccount(to)
	if err != nil {
		return err
	}
	newFromBalance := fromAccount.Balance
	if !allowsBalance(*fromAccount, newFromBalance) {
		return fmt.Errorf("%w na conta origem %s", ErrInsufficientBalance, from)
	}
	fromAccount.Balance = newFromBalance
	toAccount.Balance += amount
	if toAccount.IsBonus() {
		toAccount.Points += pointsForReceivedTransfer(amount)
	}
	if err := s.repo.Update(*fromAccount); err != nil {
		return err
	}
	return s.repo.Update(*toAccount)
}

// RenderJuros aplica a taxa de rendimento (em %) a todas as contas poupança.
func (s *AccountService) RenderJuros(taxa float64) error {
	if taxa <= 0 {
		return ErrInvalidRate
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
		return fmt.Errorf("%w: %s", ErrAccountAlreadyExists, account.Number)
	}
	return s.repo.Save(account)
}

func validateAmount(amount float64) error {
	if amount < 0 {
		return ErrNegativeAmount
	}
	if amount == 0 {
		return ErrNonPositiveAmount
	}
	return nil
}

func pointsForDeposit(amount float64) int {
	return int(amount / 100)
}

func pointsForReceivedTransfer(amount float64) int {
	return int(amount / 200)
}

func allowsBalance(account domain.Account, balance float64) bool {
	if account.Type == domain.AccountTypeSimple || account.IsBonus() {
		return balance >= negativeBalanceLimit
	}
	return balance >= 0
}
