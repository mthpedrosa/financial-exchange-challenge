# Financial Exchange Challenge

Este projeto é uma API de exchange financeira, construída em Go, que gerencia contas, balances, instrumentos e ordens. As ordens criadas são publicadas em uma fila RabbitMQ para serem processadas por outro serviço responsável pelo motor de matching.

---

## Dependências Necessárias

- **Go 1.24.8+**
- **PostgreSQL** (para persistência dos dados)
- **RabbitMQ** (para fila de ordens)
- **Docker** (opcional, para facilitar o setup dos serviços)
- **swag** (para documentação Swagger)
- Bibliotecas Go:
  - `github.com/labstack/echo/v4`
  - `github.com/jackc/pgx/v5/pgxpool`
  - `github.com/rabbitmq/amqp091-go`
  - `github.com/swaggo/echo-swagger`
  - `github.com/go-playground/validator/v10`
  - `github.com/stretchr/testify`
  - Outras listadas em `go.mod`

---

## Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis:

```
APP_ENV=development
APP_NAME=financial-exchange
LOG_LEVEL=info

DATABASE_URL=postgres://user:password@localhost:5432/financial_exchange?sslmode=disable
RABBIT_URL=amqp://guest:guest@localhost:5672/
```

Ajuste os valores conforme seu ambiente.

---

## Comandos para Rodar

1. **Instale as dependências Go:**
   ```sh
   go mod tidy
   ```

2. **Gere a documentação Swagger:**
   ```sh
   go install github.com/swaggo/swag/cmd/swag@latest
   swag init
   ```

3. **Suba o banco e o RabbitMQ (exemplo com Docker Compose):**
   ```yaml
   # docker-compose.yml
   version: '3'
   services:
     db:
       image: postgres:15
       environment:
         POSTGRES_USER: user
         POSTGRES_PASSWORD: password
         POSTGRES_DB: financial_exchange
       ports:
         - "5432:5432"
     rabbitmq:
       image: rabbitmq:3-management
       ports:
         - "5672:5672"
         - "15672:15672"
   ```
   ```sh
   docker-compose up -d
   ```

4. **Rode as migrations (se aplicável):**
   ```sh
   go run cmd/main.go migrate
   ```

5. **Inicie a aplicação:**
   ```sh
   go run cmd/main.go
   ```

6. **Acesse a documentação Swagger:**
   ```
   http://localhost:8080/swagger/index.html
   ```

---

## Integração com o Motor de Matching

Este repositório **não executa o matching das ordens**.  
Ele apenas publica as ordens criadas em uma fila RabbitMQ (`orders`).  
Outro serviço (motor de matching) deve ser responsável por consumir essa fila, processar as ordens e atualizar o status conforme necessário.

> **Importante:**  
> Certifique-se de rodar o serviço de matching em conjunto para que as ordens sejam processadas corretamente.

---

## Testes

Para rodar os testes:
```sh
go test ./...
```

---

## Observações

- O projeto segue arquitetura hexagonal (ports & adapters).
- Todas as rotas estão documentadas via Swagger.
- O código está preparado para produção, mas revise as configurações de segurança antes de expor publicamente.

---