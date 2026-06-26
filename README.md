# BitBank

## Integrantes 

> André Luiz de Sena Liberato
> Rian Abdias Balbino de Azevedo

## Tecnologias Utilizadas

> Go Lang (1.26+)
> SQLite
> Docker

---

## Como Executar a Aplicação

### 1. Localmente (Go)
Certifique-se de ter o Go instalado (versão 1.22 ou superior).

* **Modo Interativo (CLI):**
  ```bash
  go run main.go
  ```
* **Iniciar Servidor API REST:**
  ```bash
  go run main.go serve --port 8080
  ```

### 2. Com Docker
Você pode baixar e executar a imagem diretamente do Docker Hub:

```bash
docker run -d -p 8080:8080 andreliberato/bitbank:latest
```

* **Repositório no Docker Hub:** [andreliberato/bitbank](https://hub.docker.com/r/andreliberato/bitbank)

---

## Exemplos de Acesso à API REST

A API roda por padrão na porta `8080`. Abaixo estão os comandos `curl` de exemplo para interagir com a API:

### 1. Criar uma Conta
* **Tipos de conta disponíveis:** `simples`, `bonus`, `poupanca`.
* **Endpoint:** `POST /banco/conta/`
* **Exemplo:**
  ```bash
  curl -X POST http://localhost:8080/banco/conta/ \
    -H "Content-Type: application/json" \
    -d '{"number": "1234-5", "type": "poupanca", "balance": 1000.0}'
  ```

### 2. Consultar Conta (Detalhes)
* **Endpoint:** `GET /banco/conta/{id}`
* **Exemplo:**
  ```bash
  curl http://localhost:8080/banco/conta/1234-5
  ```

### 3. Consultar Saldo
* **Endpoint:** `GET /banco/conta/{id}/saldo`
* **Exemplo:**
  ```bash
  curl http://localhost:8080/banco/conta/1234-5/saldo
  ```

### 4. Creditar um Valor
* **Endpoint:** `PUT /banco/conta/{id}/credito`
* **Exemplo:**
  ```bash
  curl -X PUT http://localhost:8080/banco/conta/1234-5/credito \
    -H "Content-Type: application/json" \
    -d '{"amount": 150.0}'
  ```

### 5. Debitar um Valor
* **Endpoint:** `PUT /banco/conta/{id}/debito`
* **Exemplo:**
  ```bash
  curl -X PUT http://localhost:8080/banco/conta/1234-5/debito \
    -H "Content-Type: application/json" \
    -d '{"amount": 50.0}'
  ```

### 6. Transferência entre Contas
* **Endpoint:** `PUT /banco/conta/transferencia`
* **Exemplo:**
  ```bash
  curl -X PUT http://localhost:8080/banco/conta/transferencia \
    -H "Content-Type: application/json" \
    -d '{"from": "1234-5", "to": "6789-0", "amount": 200.0}'
  ```

### 7. Render Rendimento (Apenas para Poupança)
* **Endpoint:** `PUT /banco/conta/rendimento`
* **Exemplo:**
  ```bash
  curl -X PUT http://localhost:8080/banco/conta/rendimento \
    -H "Content-Type: application/json" \
    -d '{"rate": 0.05}'
  ```

---

## Fluxo de CI/CD

Este projeto utiliza GitHub Actions para automatizar a integração e entrega contínua.

1. **Pipeline de Integração Contínua (CI):** Executa em pushes e PRs para a branch `main`. Constrói o projeto, roda testes e gera uma tag `build-XXX`.
2. **Pipeline de Estabilização (Homologação):** Executa em PRs para as branches `stabilization/rc-*`. Executa testes, análise estática (`go vet`), cria tags `rc-*` correspondentes e gera o pacote `.zip` da release.
3. **Pipeline de Produção:** Executa em PRs para a branch `production`. Executa todos os testes e análises estáticas, gera a tag de release de produção (`rel-*`), publica o pacote `.zip` nos artefatos do workflow, constrói e publica a imagem Docker no Docker Hub.
