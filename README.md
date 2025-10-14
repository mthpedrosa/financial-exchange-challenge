# Financial Exchange Challenge 🚀

[![Go version](https://img.shields.io/badge/go-1.22+-00ADD8.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

API de uma exchange financeira construída em Go. A arquitetura gerencia contas, balances, instrumentos e ordens de compra/venda. As ordens são publicadas em uma fila **RabbitMQ** para serem processadas por um motor de matching externo.

## ✨ Features

- **Arquitetura Hexagonal:** Código organizado, testável e de fácil manutenção.
- **API RESTful:** Endpoints para gerenciar contas, instrumentos, balances e ordens.
- **Documentação com Swagger:** Interface interativa para explorar e testar a API.
- **Mensageria com RabbitMQ:** Desacoplamento para processamento assíncrono de ordens.
- **Totalmente Containerizado:** Ambiente de desenvolvimento e produção padronizado com Docker.

---

## 🚀 Rodando com Docker (Método Recomendado)

A forma mais simples e rápida de ter todo o ambiente (API, Banco de Dados e Fila) rodando.

### Pré-requisitos

- [Docker](https://www.docker.com/get-started/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Passo 1: Clone o Repositório

```sh
git clone <url-do-seu-repositorio>
cd financial-exchange-challenge
````

### Passo 2: Crie o Arquivo de Ambiente (`.env`)

Crie um arquivo chamado `.env` na raiz do projeto. **Copie e cole** o conteúdo abaixo. Note que usamos os nomes dos serviços (`db`, `rabbitmq`) em vez de `localhost`.

```env
# .env para o ambiente Docker

# Configurações da Aplicação
APP_ENV=development
PORT=8080
LOG_LEVEL=debug
APP_NAME="Exchange API"

# Segredo para JWT
JWT_SECRET="uma-chave-secreta-forte-e-aleatoria-aqui"

# Conexões com os Serviços do Docker Compose
DATABASE_URL=postgresql://user:password@db:5432/financial_exchange?sslmode=disable
RABBITMQ_URL="amqp://guest:guest@rabbitmq:5672/"
```

### Passo 3: Suba os Contêineres

Este único comando irá construir a imagem da sua API, baixar as imagens do Postgres e RabbitMQ, e iniciar todos os serviços em segundo plano.

```sh
docker-compose up --build -d
```

  * `--build`: Constrói a imagem da sua API. Use sempre que alterar o código Go.
  * `-d`: Roda os contêineres em modo "detached" (segundo plano).

### Pronto\!

O ambiente está no ar.

  * **API disponível em:** `http://localhost:8080`
  * **Documentação Swagger:** `http://localhost:8080/swagger/index.html`
  * **Painel do RabbitMQ:** `http://localhost:15672` (usuário: `guest`, senha: `guest`)

-----

## 🛠️ Rodando Localmente (Para Desenvolvimento)

Use este método se você preferir rodar a aplicação Go diretamente na sua máquina, mas usando Docker para as dependências.

### Pré-requisitos

  - Go (versão 1.22+)
  - Docker e Docker Compose
  - [Swag CLI](https://github.com/swaggo/swag) (`go install github.com/swaggo/swag/cmd/swag@latest`)

### Passo 1: Suba as Dependências

Inicie apenas o banco de dados e o RabbitMQ com Docker Compose.

```sh
docker-compose up -d db rabbitmq
```

### Passo 2: Crie o Arquivo de Ambiente (`.env`)

Crie o arquivo `.env`, mas desta vez, use `localhost` para as conexões, pois sua API estará rodando fora do Docker.

```env
# .env para o ambiente Local

# ... (demais variáveis como APP_ENV, PORT, etc.)

# Conexões com os Serviços (via localhost)
DATABASE_URL=postgresql://user:password@localhost:5432/financial_exchange?sslmode=disable
RABBITMQ_URL="amqp://guest:guest@localhost:5672/"
```

### Passo 3: Instale as Dependências e Gere a Documentação

```sh
go mod tidy
swag init -g cmd/main.go -parseDependency
```

### Passo 4: Inicie a Aplicação Go

```sh
go run ./cmd/main.go
```

A API estará rodando e conectada aos serviços do Docker.

-----

## 🧪 Testes

Para rodar todos os testes unitários e de integração:

```sh
go test ./...
```

-----

## 🏛️ Arquitetura

O projeto utiliza uma abordagem de **Arquitetura Hexagonal (Ports and Adapters)** para separar as regras de negócio da infraestrutura. Isso resulta em um código mais limpo, desacoplado e fácil de testar.

  - **`internal/`**: Contém o núcleo da aplicação (domínio, casos de uso) e as implementações dos adaptadores (handlers de API, repositórios de banco de dados).
  - **`cmd/`**: Ponto de entrada da aplicação, onde tudo é inicializado e conectado.

-----

---
## 👨‍💻 Desenvolvido por

[<img src="https://upload.wikimedia.org/wikipedia/commons/thumb/8/81/LinkedIn_icon.svg/1024px-LinkedIn_icon.svg.png" width="100px;" alt="Matheus Pedrosa"/><br><sub><b>Matheus Pedrosa</b></sub>](https://www.linkedin.com/in/matheus-pedrosa-custodio/)

```