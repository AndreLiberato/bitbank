package service

import "errors"

// Erros de negócio do banco. Funcionam como uma hierarquia de "exceções":
// a camada de serviço sempre retorna (ou encapsula com %w) um destes erros,
// permitindo que as camadas superiores (REST, CLI) reajam de forma tipada
// via errors.Is, sem inspecionar mensagens.
var (
	// ErrNegativeAmount indica que um valor monetário informado é negativo.
	ErrNegativeAmount = errors.New("valor não pode ser negativo")
	// ErrNonPositiveAmount indica que um valor deveria ser maior que zero.
	ErrNonPositiveAmount = errors.New("valor deve ser maior que zero")
	// ErrNegativeInitialBalance indica saldo inicial negativo no cadastro.
	ErrNegativeInitialBalance = errors.New("saldo inicial não pode ser negativo")
	// ErrInsufficientBalance indica que a operação deixaria o saldo abaixo do permitido.
	ErrInsufficientBalance = errors.New("saldo insuficiente")
	// ErrAccountNotFound indica que a conta informada não existe.
	ErrAccountNotFound = errors.New("conta não encontrada")
	// ErrAccountAlreadyExists indica tentativa de cadastrar conta já existente.
	ErrAccountAlreadyExists = errors.New("conta já existe")
	// ErrInvalidAccountType indica um tipo de conta desconhecido no cadastro.
	ErrInvalidAccountType = errors.New("tipo de conta inválido")
	// ErrInvalidRate indica taxa de rendimento inválida.
	ErrInvalidRate = errors.New("taxa deve ser maior que zero")
)
