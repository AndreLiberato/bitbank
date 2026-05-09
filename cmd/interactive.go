package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/AndreLiberato/bitbank/internal/domain"
	"github.com/AndreLiberato/bitbank/internal/service"
	"github.com/charmbracelet/huh"
)

const (
	opCadastrar     = "cadastrar"
	opSaldo         = "saldo"
	opCredito       = "credito"
	opDebito        = "debito"
	opTransferencia = "transferencia"
	opSair          = "sair"
)

func RunInteractive(svc *service.AccountService) {
	for {
		op := selecionarOperacao()
		if op == opSair {
			fmt.Println("Até logo!")
			return
		}
		executarOperacao(op, svc)
		aguardarEnter()
	}
}

func selecionarOperacao() string {
	var op string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("BitBank — escolha uma operação").
				Options(
					huh.NewOption("Cadastrar Conta", opCadastrar),
					huh.NewOption("Consultar Saldo", opSaldo),
					huh.NewOption("Crédito", opCredito),
					huh.NewOption("Débito", opDebito),
					huh.NewOption("Transferência", opTransferencia),
					huh.NewOption("Sair", opSair),
				).
				Value(&op),
		),
	)
	if err := form.Run(); err != nil {
		os.Exit(0)
	}
	return op
}

func executarOperacao(op string, svc *service.AccountService) {
	switch op {
	case opCadastrar:
		cadastrarConta(svc)
	case opSaldo:
		consultarSaldo(svc)
	case opCredito:
		credito(svc)
	case opDebito:
		debito(svc)
	case opTransferencia:
		transferencia(svc)
	}
}

func cadastrarConta(svc *service.AccountService) {
	var numero, tipo string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Cadastrar Conta").
				Description("Informe o número da nova conta").
				Value(&numero).
				Validate(naoVazio),
			huh.NewSelect[string]().
				Title("Tipo de conta").
				Options(
					huh.NewOption("Conta Simples", domain.AccountTypeSimple),
					huh.NewOption("Conta Bônus", domain.AccountTypeBonus),
				).
				Value(&tipo),
		),
	)
	if err := form.Run(); err != nil {
		return
	}
	if tipo == domain.AccountTypeBonus {
		if err := svc.CreateBonusAccount(numero); err != nil {
			printErro(err)
			return
		}
		printSucesso(fmt.Sprintf("Conta bônus %s criada com saldo R$ 0,00 e 10 pontos.", numero))
		return
	}
	if err := svc.CreateAccount(numero); err != nil {
		printErro(err)
		return
	}
	printSucesso(fmt.Sprintf("Conta %s criada com saldo R$ 0,00.", numero))
}

func consultarSaldo(svc *service.AccountService) {
	var numero string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Consultar Saldo").
				Description("Informe o número da conta").
				Value(&numero).
				Validate(naoVazio),
		),
	)
	if err := form.Run(); err != nil {
		return
	}
	account, err := svc.GetAccount(numero)
	if err != nil {
		printErro(err)
		return
	}
	msg := fmt.Sprintf("Saldo da conta %s: R$ %.2f", numero, account.Balance)
	if account.IsBonus() {
		msg = fmt.Sprintf("%s | Pontos: %d", msg, account.Points)
	}
	printSucesso(msg)
}

func credito(svc *service.AccountService) {
	var numero, valorStr string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Crédito").
				Description("Número da conta").
				Value(&numero).
				Validate(naoVazio),
			huh.NewInput().
				Description("Valor").
				Value(&valorStr).
				Validate(validarValor),
		),
	)
	if err := form.Run(); err != nil {
		return
	}
	if err := svc.Credit(numero, parseValor(valorStr)); err != nil {
		printErro(err)
		return
	}
	printSucesso(fmt.Sprintf("Crédito de R$ %.2f realizado na conta %s.", parseValor(valorStr), numero))
}

func debito(svc *service.AccountService) {
	var numero, valorStr string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Débito").
				Description("Número da conta").
				Value(&numero).
				Validate(naoVazio),
			huh.NewInput().
				Description("Valor").
				Value(&valorStr).
				Validate(validarValor),
		),
	)
	if err := form.Run(); err != nil {
		return
	}
	if err := svc.Debit(numero, parseValor(valorStr)); err != nil {
		printErro(err)
		return
	}
	printSucesso(fmt.Sprintf("Débito de R$ %.2f realizado na conta %s.", parseValor(valorStr), numero))
}

func transferencia(svc *service.AccountService) {
	var origem, destino, valorStr string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Transferência").
				Description("Conta origem").
				Value(&origem).
				Validate(naoVazio),
			huh.NewInput().
				Description("Conta destino").
				Value(&destino).
				Validate(naoVazio),
			huh.NewInput().
				Description("Valor").
				Value(&valorStr).
				Validate(validarValor),
		),
	)
	if err := form.Run(); err != nil {
		return
	}
	valor := parseValor(valorStr)
	if err := svc.Transfer(origem, destino, valor); err != nil {
		printErro(err)
		return
	}
	printSucesso(fmt.Sprintf("Transferência de R$ %.2f da conta %s para %s realizada.", valor, origem, destino))
}

func aguardarEnter() {
	fmt.Print("\nPressione Enter para continuar...")
	fmt.Scanln()
}

func printSucesso(msg string) {
	fmt.Printf("\n✓ %s\n", msg)
}

func printErro(err error) {
	fmt.Printf("\n✗ Erro: %s\n", err)
}

func naoVazio(s string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("campo obrigatório")
	}
	return nil
}

func validarValor(s string) error {
	v, err := strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
	if err != nil {
		return fmt.Errorf("valor inválido")
	}
	if v <= 0 {
		return fmt.Errorf("valor deve ser maior que zero")
	}
	return nil
}

func parseValor(s string) float64 {
	v, _ := strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
	return v
}
