package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/AndreLiberato/bitbank/internal/service"
	"github.com/charmbracelet/huh"
)

const (
	opCadastrar = "cadastrar"
	opSair      = "sair"
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
	}
}

func cadastrarConta(svc *service.AccountService) {
	var numero string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Cadastrar Conta").
				Description("Informe o número da nova conta").
				Value(&numero).
				Validate(naoVazio),
		),
	)
	if err := form.Run(); err != nil {
		return
	}
	if err := svc.CreateAccount(numero); err != nil {
		printErro(err)
		return
	}
	printSucesso(fmt.Sprintf("Conta %s criada com saldo R$ 0,00.", numero))
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
