package api

import (
	"encoding/json"
	"net/http"

	"github.com/AndreLiberato/bitbank/internal/service"
)

// AccountHandler é o "controller" REST das contas. Não contém regra de negócio:
// apenas traduz HTTP <-> chamadas à fachada (service.AccountService).
type AccountHandler struct {
	svc *service.AccountService
}

// NewAccountHandler cria o controller a partir da fachada de negócios.
func NewAccountHandler(svc *service.AccountService) *AccountHandler {
	return &AccountHandler{svc: svc}
}

// Create — POST /banco/conta/
func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountRequest
	if !decode(w, r, &req) {
		return
	}
	if err := h.svc.OpenAccount(req.Number, req.Type, req.Balance); err != nil {
		writeError(w, err)
		return
	}
	account, err := h.svc.GetAccount(req.Number)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toAccountResponse(account))
}

// Get — GET /banco/conta/{id}
func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	account, err := h.svc.GetAccount(id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toAccountResponse(account))
}

// Balance — GET /banco/conta/{id}/saldo
func (h *AccountHandler) Balance(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	balance, err := h.svc.GetBalance(id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, BalanceResponse{Number: id, Balance: balance})
}

// Credit — PUT /banco/conta/{id}/credito
func (h *AccountHandler) Credit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req AmountRequest
	if !decode(w, r, &req) {
		return
	}
	if err := h.svc.Credit(id, req.Amount); err != nil {
		writeError(w, err)
		return
	}
	h.respondWithAccount(w, id)
}

// Debit — PUT /banco/conta/{id}/debito
func (h *AccountHandler) Debit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req AmountRequest
	if !decode(w, r, &req) {
		return
	}
	if err := h.svc.Debit(id, req.Amount); err != nil {
		writeError(w, err)
		return
	}
	h.respondWithAccount(w, id)
}

// Transfer — PUT /banco/conta/transferencia
func (h *AccountHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req TransferRequest
	if !decode(w, r, &req) {
		return
	}
	if err := h.svc.Transfer(req.From.String(), req.To.String(), req.Amount); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, MessageResponse{
		Message: "transferência realizada com sucesso",
	})
}

// RenderJuros — PUT /banco/conta/rendimento
func (h *AccountHandler) RenderJuros(w http.ResponseWriter, r *http.Request) {
	var req RenderJurosRequest
	if !decode(w, r, &req) {
		return
	}
	if err := h.svc.RenderJuros(req.Rate); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, MessageResponse{
		Message: "rendimento aplicado às contas poupança",
	})
}

// respondWithAccount consulta a conta e devolve seu estado atualizado.
func (h *AccountHandler) respondWithAccount(w http.ResponseWriter, id string) {
	account, err := h.svc.GetAccount(id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toAccountResponse(account))
}

// ---- helpers de transporte ----

// decode lê e valida o corpo JSON da requisição. Em caso de erro, já responde
// 400 e retorna false, sinalizando ao handler para interromper.
func decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "corpo da requisição inválido"})
		return false
	}
	return true
}

func writeError(w http.ResponseWriter, err error) {
	writeJSON(w, httpStatusFor(err), ErrorResponse{Error: err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
