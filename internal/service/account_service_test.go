package service

import (
	"errors"
	"testing"

	"github.com/AndreLiberato/bitbank/internal/domain"
	"github.com/AndreLiberato/bitbank/internal/repository"
)

// newService cria uma fachada apoiada em um repositório em memória, isolando
// cada teste do SQLite e exercitando apenas as regras de negócio.
func newService() *AccountService {
	return NewAccountService(repository.NewMemoryRepository())
}

// ---------------------------------------------------------------------------
// Cadastrar Conta — um teste para cada tipo de conta
// ---------------------------------------------------------------------------

func TestCreateAccount_Simple(t *testing.T) {
	svc := newService()

	if err := svc.CreateAccount("1", 100); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	acc, err := svc.GetAccount("1")
	if err != nil {
		t.Fatalf("erro ao consultar: %v", err)
	}
	if acc.Type != domain.AccountTypeSimple {
		t.Errorf("tipo esperado %q, obtido %q", domain.AccountTypeSimple, acc.Type)
	}
	if acc.Balance != 100 {
		t.Errorf("saldo esperado 100, obtido %v", acc.Balance)
	}
}

func TestCreateAccount_Bonus(t *testing.T) {
	svc := newService()

	if err := svc.CreateBonusAccount("2"); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	acc, _ := svc.GetAccount("2")
	if acc.Type != domain.AccountTypeBonus {
		t.Errorf("tipo esperado %q, obtido %q", domain.AccountTypeBonus, acc.Type)
	}
	if acc.Balance != 0 {
		t.Errorf("saldo inicial esperado 0, obtido %v", acc.Balance)
	}
	if acc.Points != 10 {
		t.Errorf("conta bônus deve iniciar com 10 pontos, obtido %d", acc.Points)
	}
}

func TestCreateAccount_Savings(t *testing.T) {
	svc := newService()

	if err := svc.CreateSavingsAccount("3", 250); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	acc, _ := svc.GetAccount("3")
	if acc.Type != domain.AccountTypeSavings {
		t.Errorf("tipo esperado %q, obtido %q", domain.AccountTypeSavings, acc.Type)
	}
	if acc.Balance != 250 {
		t.Errorf("saldo esperado 250, obtido %v", acc.Balance)
	}
}

func TestCreateAccount_DuplicateFails(t *testing.T) {
	svc := newService()
	_ = svc.CreateAccount("1", 0)

	err := svc.CreateAccount("1", 0)
	if !errors.Is(err, ErrAccountAlreadyExists) {
		t.Errorf("esperado ErrAccountAlreadyExists, obtido %v", err)
	}
}

func TestOpenAccount_DispatchesByType(t *testing.T) {
	svc := newService()

	cases := []struct {
		number   string
		typ      string
		wantType string
	}{
		{"10", domain.AccountTypeSimple, domain.AccountTypeSimple},
		{"20", domain.AccountTypeBonus, domain.AccountTypeBonus},
		{"30", domain.AccountTypeSavings, domain.AccountTypeSavings},
	}
	for _, c := range cases {
		if err := svc.OpenAccount(c.number, c.typ, 50); err != nil {
			t.Fatalf("OpenAccount(%s): erro inesperado: %v", c.typ, err)
		}
		acc, _ := svc.GetAccount(c.number)
		if acc.Type != c.wantType {
			t.Errorf("tipo esperado %q, obtido %q", c.wantType, acc.Type)
		}
	}

	if err := svc.OpenAccount("99", "investimento", 0); !errors.Is(err, ErrInvalidAccountType) {
		t.Errorf("esperado ErrInvalidAccountType, obtido %v", err)
	}
}

// ---------------------------------------------------------------------------
// Consultar Conta — um teste para cada tipo de conta
// ---------------------------------------------------------------------------

func TestGetAccount_Simple(t *testing.T) {
	svc := newService()
	_ = svc.CreateAccount("1", 100)

	acc, err := svc.GetAccount("1")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if acc.Number != "1" || acc.Type != domain.AccountTypeSimple || acc.Balance != 100 {
		t.Errorf("dados inesperados: %+v", acc)
	}
}

func TestGetAccount_Bonus(t *testing.T) {
	svc := newService()
	_ = svc.CreateBonusAccount("2")

	acc, err := svc.GetAccount("2")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if acc.Type != domain.AccountTypeBonus || acc.Points != 10 {
		t.Errorf("dados inesperados: %+v", acc)
	}
}

func TestGetAccount_Savings(t *testing.T) {
	svc := newService()
	_ = svc.CreateSavingsAccount("3", 300)

	acc, err := svc.GetAccount("3")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if acc.Type != domain.AccountTypeSavings || acc.Balance != 300 {
		t.Errorf("dados inesperados: %+v", acc)
	}
}

func TestGetAccount_NotFound(t *testing.T) {
	svc := newService()

	_, err := svc.GetAccount("inexistente")
	if !errors.Is(err, ErrAccountNotFound) {
		t.Errorf("esperado ErrAccountNotFound, obtido %v", err)
	}
}

// ---------------------------------------------------------------------------
// Consultar Saldo
// ---------------------------------------------------------------------------

func TestGetBalance(t *testing.T) {
	svc := newService()
	_ = svc.CreateAccount("1", 175.50)

	balance, err := svc.GetBalance("1")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if balance != 175.50 {
		t.Errorf("saldo esperado 175.50, obtido %v", balance)
	}
}

func TestGetBalance_NotFound(t *testing.T) {
	svc := newService()

	if _, err := svc.GetBalance("x"); !errors.Is(err, ErrAccountNotFound) {
		t.Errorf("esperado ErrAccountNotFound, obtido %v", err)
	}
}

// ---------------------------------------------------------------------------
// Crédito
// ---------------------------------------------------------------------------

func TestCredit_Normal(t *testing.T) {
	svc := newService()
	_ = svc.CreateAccount("1", 100)

	if err := svc.Credit("1", 50); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	balance, _ := svc.GetBalance("1")
	if balance != 150 {
		t.Errorf("saldo esperado 150, obtido %v", balance)
	}
}

func TestCredit_NegativeAmountRejected(t *testing.T) {
	svc := newService()
	_ = svc.CreateAccount("1", 100)

	if err := svc.Credit("1", -10); !errors.Is(err, ErrNegativeAmount) {
		t.Errorf("esperado ErrNegativeAmount, obtido %v", err)
	}

	balance, _ := svc.GetBalance("1")
	if balance != 100 {
		t.Errorf("saldo não deveria mudar, obtido %v", balance)
	}
}

func TestCredit_BonusAccumulatesPoints(t *testing.T) {
	svc := newService()
	_ = svc.CreateBonusAccount("2") // inicia com 10 pontos

	if err := svc.Credit("2", 300); err != nil { // +3 pontos (1 a cada 100)
		t.Fatalf("erro inesperado: %v", err)
	}

	acc, _ := svc.GetAccount("2")
	if acc.Balance != 300 {
		t.Errorf("saldo esperado 300, obtido %v", acc.Balance)
	}
	if acc.Points != 13 {
		t.Errorf("pontos esperados 13 (10 + 3), obtido %d", acc.Points)
	}
}

// ---------------------------------------------------------------------------
// Débito
// ---------------------------------------------------------------------------

func TestDebit_Normal(t *testing.T) {
	svc := newService()
	_ = svc.CreateAccount("1", 100)

	if err := svc.Debit("1", 40); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	balance, _ := svc.GetBalance("1")
	if balance != 60 {
		t.Errorf("saldo esperado 60, obtido %v", balance)
	}
}

func TestDebit_NegativeAmountRejected(t *testing.T) {
	svc := newService()
	_ = svc.CreateAccount("1", 100)

	if err := svc.Debit("1", -5); !errors.Is(err, ErrNegativeAmount) {
		t.Errorf("esperado ErrNegativeAmount, obtido %v", err)
	}
}

func TestDebit_PreventsNegativeBalance(t *testing.T) {
	svc := newService()
	// Conta poupança não admite saldo negativo.
	_ = svc.CreateSavingsAccount("3", 50)

	if err := svc.Debit("3", 100); !errors.Is(err, ErrInsufficientBalance) {
		t.Errorf("esperado ErrInsufficientBalance, obtido %v", err)
	}

	balance, _ := svc.GetBalance("3")
	if balance != 50 {
		t.Errorf("saldo não deveria mudar, obtido %v", balance)
	}
}

// ---------------------------------------------------------------------------
// Transferência
// ---------------------------------------------------------------------------

func TestTransfer_NegativeAmountRejected(t *testing.T) {
	svc := newService()
	_ = svc.CreateAccount("1", 100)
	_ = svc.CreateAccount("2", 0)

	if err := svc.Transfer("1", "2", -50); !errors.Is(err, ErrNegativeAmount) {
		t.Errorf("esperado ErrNegativeAmount, obtido %v", err)
	}
}

func TestTransfer_PreventsNegativeBalance(t *testing.T) {
	svc := newService()
	_ = svc.CreateSavingsAccount("3", 50) // origem não pode ficar negativa
	_ = svc.CreateAccount("1", 0)

	if err := svc.Transfer("3", "1", 100); !errors.Is(err, ErrInsufficientBalance) {
		t.Errorf("esperado ErrInsufficientBalance, obtido %v", err)
	}

	origem, _ := svc.GetBalance("3")
	destino, _ := svc.GetBalance("1")
	if origem != 50 || destino != 0 {
		t.Errorf("saldos não deveriam mudar: origem=%v destino=%v", origem, destino)
	}
}

func TestTransfer_BonusDestinationAccumulatesPoints(t *testing.T) {
	svc := newService()
	_ = svc.CreateAccount("1", 1000)
	_ = svc.CreateBonusAccount("2") // inicia com 10 pontos

	if err := svc.Transfer("1", "2", 400); err != nil { // +2 pontos (1 a cada 200)
		t.Fatalf("erro inesperado: %v", err)
	}

	destino, _ := svc.GetAccount("2")
	if destino.Balance != 400 {
		t.Errorf("saldo destino esperado 400, obtido %v", destino.Balance)
	}
	if destino.Points != 12 {
		t.Errorf("pontos esperados 12 (10 + 2), obtido %d", destino.Points)
	}

	origem, _ := svc.GetBalance("1")
	if origem != 600 {
		t.Errorf("saldo origem esperado 600, obtido %v", origem)
	}
}

// ---------------------------------------------------------------------------
// Render Juros
// ---------------------------------------------------------------------------

func TestRenderJuros_AppliesToAllSavingsAccounts(t *testing.T) {
	svc := newService()
	_ = svc.CreateSavingsAccount("p1", 100)
	_ = svc.CreateSavingsAccount("p2", 200)
	_ = svc.CreateAccount("s1", 100)   // simples — não rende
	_ = svc.CreateBonusAccount("b1")   // bônus — não rende

	if err := svc.RenderJuros(10); err != nil { // 10%
		t.Fatalf("erro inesperado: %v", err)
	}

	if b, _ := svc.GetBalance("p1"); b != 110 {
		t.Errorf("poupança p1 esperada 110, obtido %v", b)
	}
	if b, _ := svc.GetBalance("p2"); b != 220 {
		t.Errorf("poupança p2 esperada 220, obtido %v", b)
	}
	if b, _ := svc.GetBalance("s1"); b != 100 {
		t.Errorf("conta simples não deveria render, obtido %v", b)
	}
	if b, _ := svc.GetBalance("b1"); b != 0 {
		t.Errorf("conta bônus não deveria render, obtido %v", b)
	}
}

func TestRenderJuros_InvalidRateRejected(t *testing.T) {
	svc := newService()

	if err := svc.RenderJuros(0); !errors.Is(err, ErrInvalidRate) {
		t.Errorf("esperado ErrInvalidRate, obtido %v", err)
	}
}
