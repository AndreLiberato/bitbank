package api

import (
	"errors"
	"net/http"

	"github.com/AndreLiberato/bitbank/internal/service"
)

// httpStatusFor traduz os erros de negócio da camada de serviço em códigos HTTP.
// Centralizar esse mapeamento aqui mantém os handlers limpos e desacopla o
// transporte (HTTP) das regras de negócio.
func httpStatusFor(err error) int {
	switch {
	case errors.Is(err, service.ErrAccountNotFound):
		return http.StatusNotFound
	case errors.Is(err, service.ErrAccountAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, service.ErrNegativeAmount),
		errors.Is(err, service.ErrNonPositiveAmount),
		errors.Is(err, service.ErrNegativeInitialBalance),
		errors.Is(err, service.ErrInsufficientBalance),
		errors.Is(err, service.ErrInvalidAccountType),
		errors.Is(err, service.ErrInvalidRate):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
