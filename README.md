# bitbank

Sistema bancário de linha de comando (CLI) desenvolvido em Go com interface interativa.

## Integrantes

> André Luiz de Sena Liberato

> Rian Abdias Balbino de Azevedo

## Tecnologias Utilizadas

- [Go](https://golang.org/) 1.26+
- [Cobra](https://github.com/spf13/cobra) — framework CLI
- [charmbracelet/huh](https://github.com/charmbracelet/huh) — interface interativa de terminal
- [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) — banco de dados SQLite (puro Go, sem CGO)

## Funcionalidades

| Operação | Descrição |
|----------|-----------|
| Cadastrar Conta | Cria uma conta com número informado e saldo inicial zero |
| Consultar Saldo | Exibe o saldo atual de uma conta |
| Crédito | Adiciona valor ao saldo de uma conta |
| Débito | Subtrai valor do saldo de uma conta (saldo negativo permitido) |
| Transferência | Move valor entre duas contas |

## Arquitetura

O projeto segue arquitetura em camadas com separação de responsabilidades:

```
cmd/                        → Camada de apresentação (CLI)
  root.go                   → Ponto de entrada Cobra, inicializa dependências
  interactive.go            → Loop interativo, menus e formulários (huh)

internal/
  domain/
    account.go              → Entidade Account {Number, Balance}

  repository/
    account_repository.go   → Interface + implementação SQLite

  service/
    account_service.go      → Regras de negócio

main.go                     → Entrypoint
```

**Fluxo de dependências:** `cmd → service → repository → domain`

## Banco de Dados

O banco SQLite é criado automaticamente na primeira execução em:

```
~/.bitbank/bitbank.db
```

O arquivo não é versionado no repositório.

## Como Executar

**Pré-requisito:** Go 1.21 ou superior instalado.

```bash
# Clonar o repositório
git clone https://github.com/AndreLiberato/bitbank.git
cd bitbank

# Executar diretamente
go run main.go

# Ou compilar e executar
go build -o bitbank .
./bitbank
```

## Regras de Negócio

- Número de contas ilimitado
- Contas podem ter saldo negativo
- Conta possui apenas número e saldo
- Valores de crédito, débito e transferência devem ser maiores que zero
- Não é possível cadastrar duas contas com o mesmo número
