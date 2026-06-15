package api

import (
	"encoding/json"

	"github.com/AndreLiberato/bitbank/internal/domain"
)

// DTOs de entrada (request bodies) e saída (responses) da API REST.
// Mantê-los separados do domínio evita acoplar o contrato HTTP às entidades.

// CreateAccountRequest representa o corpo do cadastro de conta.
type CreateAccountRequest struct {
	Number  string  `json:"number"`
	Type    string  `json:"type"`
	Balance float64 `json:"balance"`
}

// AmountRequest representa o corpo de operações de crédito e débito.
type AmountRequest struct {
	Amount float64 `json:"amount"`
}

// TransferRequest representa o corpo de uma transferência.
// from/to aceitam números (conforme especificação) e são convertidos para string.
type TransferRequest struct {
	From   json.Number `json:"from"`
	To     json.Number `json:"to"`
	Amount float64     `json:"amount"`
}

// RenderJurosRequest representa o corpo do rendimento (taxa em %).
type RenderJurosRequest struct {
	Rate float64 `json:"rate"`
}

// AccountResponse expõe os dados de uma conta: Tipo, Número, Saldo e Bônus
// (este último apenas quando a conta é do tipo bônus).
type AccountResponse struct {
	Type    string  `json:"type"`
	Number  string  `json:"number"`
	Balance float64 `json:"balance"`
	Bonus   *int    `json:"bonus,omitempty"`
}

// BalanceResponse expõe apenas o saldo de uma conta.
type BalanceResponse struct {
	Number  string  `json:"number"`
	Balance float64 `json:"balance"`
}

// MessageResponse é uma resposta genérica de sucesso.
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse é o corpo padrão de erro da API.
type ErrorResponse struct {
	Error string `json:"error"`
}

// toAccountResponse converte a entidade de domínio no DTO de saída.
func toAccountResponse(a *domain.Account) AccountResponse {
	resp := AccountResponse{
		Type:    a.Type,
		Number:  a.Number,
		Balance: a.Balance,
	}
	if a.IsBonus() {
		bonus := a.Points
		resp.Bonus = &bonus
	}
	return resp
}
