# 📦 User Orders API

API REST para gerenciamento de usuários e pedidos, desenvolvida em **Go** com **Gin**, **GORM** e **PostgreSQL**.

---

## 📋 Índice

- [Sobre o Projeto](#-sobre-o-projeto)
- [Tecnologias](#-tecnologias)
- [Arquitetura](#-arquitetura)
- [Estrutura do Projeto](#-estrutura-do-projeto)
- [Pré-requisitos](#-pré-requisitos)
- [Instalação e Execução](#-instalação-e-execução)
- [Variáveis de Ambiente](#-variáveis-de-ambiente)
- [Endpoints](#-endpoints)
- [Exemplos de Requisição](#-exemplos-de-requisição)
- [Tratamento de Erros](#-tratamento-de-erros)
- [Decisões de Arquitetura](#-decisões-de-arquitetura)

---

## 📌 Sobre o Projeto

Sistema backend para controle de **usuários** e seus **pedidos**. Cada usuário pode ter múltiplos pedidos com status rastreável (`pending`, `paid`, `canceled`). A API segue princípios REST com respostas padronizadas em JSON, paginação em listagens e validação de entrada em todas as rotas.

---

## 🛠 Tecnologias

| Camada | Tecnologia | Versão |
|---|---|---|
| Linguagem | Go (Golang) | 1.25 |
| Framework HTTP | Gin | v1.9.1 |
| ORM | GORM | v1.25.7 |
| Banco de Dados | PostgreSQL | 16 |
| Driver PostgreSQL | pgx via GORM | v1.5.6 |
| Logger | Uber Zap | v1.26.0 |
| Variáveis de Ambiente | godotenv | v1.5.1 |
| Autenticação | golang-jwt/jwt | v5.2.0 |
| Container | Docker + Docker Compose | - |

---

## 🏗 Arquitetura

O projeto segue o padrão de **arquitetura em camadas**, com separação clara de responsabilidades:

```
HTTP Request
     │
     ▼
┌─────────────┐
│   Handler   │  ← Valida input, chama Service, retorna JSON
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Service   │  ← Regras de negócio, orquestra chamadas
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Repository  │  ← Acesso ao banco via GORM
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ PostgreSQL  │  ← Persistência
└─────────────┘
```

**Princípios aplicados:**

- **Inversão de dependência** — cada camada depende de interfaces, não de implementações concretas
- **DTOs** — separação entre entidade de banco e contrato HTTP
- **Erros tipados** — `AppError` com código HTTP semântico
- **Middleware** — logging estruturado com Zap e recovery de panics

---

## 📁 Estrutura do Projeto

```
user-orders-api/
├── cmd/
│   └── main.go                  # Ponto de entrada — wiring de dependências
├── internal/
│   ├── config/
│   │   └── config.go            # Conexão com banco + AutoMigrate
│   ├── model/
│   │   ├── user.go              # Entidade User (mapeamento GORM)
│   │   └── order.go             # Entidade Order + OrderStatus
│   ├── dto/
│   │   ├── user_dto.go          # Request/Response de usuário
│   │   └── order_dto.go         # Request/Response de pedido + PaginationParams
│   ├── repository/
│   │   ├── user_repository.go   # Interface + impl GORM para User
│   │   └── order_repository.go  # Interface + impl GORM para Order
│   ├── service/
│   │   ├── user_service.go      # Regras de negócio de usuário
│   │   └── order_service.go     # Regras de negócio de pedido
│   └── handler/
│       ├── user_handler.go      # Controllers HTTP de usuário
│       └── order_handler.go     # Controllers HTTP de pedido
├── pkg/
│   ├── errors/
│   │   └── app_error.go         # Tipo AppError padronizado
│   └── middleware/
│       └── logger.go            # Middleware de logging com Zap
├── .env                         # Variáveis de ambiente (não versionado)
├── .env.example                 # Template de variáveis (versionado)
├── .gitignore
├── docker-compose.yml           # PostgreSQL via Docker
├── go.mod
├── go.sum
└── README.md
```

---

## ✅ Pré-requisitos

- [Go 1.21+](https://go.dev/dl/)
- [Docker](https://www.docker.com/) e Docker Compose
- `curl` + [`jq`](https://stedolan.github.io/jq/) para testes via terminal (opcional)

---

## 🚀 Instalação e Execução

### 1. Clone o repositório

```bash
git clone https://github.com/alexanderbs3/user-orders-api.git
cd user-orders-api
```

### 2. Configure as variáveis de ambiente

```bash
cp .env.example .env
# Edite o .env com suas credenciais se necessário
```

### 3. Suba o banco de dados

```bash
docker-compose up -d
```

Aguarde o healthcheck do container:

```bash
docker-compose ps
# STATUS deve ser: healthy
```

### 4. Instale as dependências

```bash
go mod tidy
```

### 5. Execute a aplicação

```bash
go run cmd/main.go
```

A API estará disponível em: `http://localhost:8080`

### Build para produção

```bash
go build -o bin/api cmd/main.go
./bin/api
```

---

## 🔧 Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto com base no `.env.example`:

```env
# Servidor
APP_PORT=8080
APP_ENV=development        # development | production

# Banco de Dados
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=user_orders_db
DB_SSLMODE=disable

# JWT (opcional)
JWT_SECRET=seu_segredo_aqui_minimo_32_caracteres
JWT_EXPIRATION_HOURS=24
```

> ⚠️ Nunca versione o arquivo `.env` com credenciais reais. Use `.env.example` para documentar as variáveis necessárias.

---

## 📡 Endpoints

### Usuários

| Método | Rota | Descrição | Status de sucesso |
|---|---|---|---|
| `POST` | `/api/v1/users` | Criar usuário | `201 Created` |
| `GET` | `/api/v1/users` | Listar usuários (paginado) | `200 OK` |
| `GET` | `/api/v1/users/:id` | Buscar usuário por ID | `200 OK` |
| `PUT` | `/api/v1/users/:id` | Atualizar usuário | `200 OK` |
| `DELETE` | `/api/v1/users/:id` | Deletar usuário | `204 No Content` |
| `GET` | `/api/v1/users/:id/orders` | Listar pedidos do usuário | `200 OK` |

### Pedidos

| Método | Rota | Descrição | Status de sucesso |
|---|---|---|---|
| `POST` | `/api/v1/orders` | Criar pedido | `201 Created` |
| `GET` | `/api/v1/orders` | Listar pedidos (paginado) | `200 OK` |
| `GET` | `/api/v1/orders/:id` | Buscar pedido por ID | `200 OK` |
| `DELETE` | `/api/v1/orders/:id` | Deletar pedido | `204 No Content` |

### Parâmetros de paginação

Todas as rotas de listagem aceitam query params:

| Parâmetro | Tipo | Padrão | Descrição |
|---|---|---|---|
| `page` | `int` | `1` | Número da página |
| `limit` | `int` | `10` | Itens por página (máx: 100) |

**Exemplo:** `GET /api/v1/users?page=2&limit=5`

---

## 📨 Exemplos de Requisição

### Criar usuário

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alexander",
    "email": "alexander@email.com"
  }'
```

**Resposta `201`:**
```json
{
  "id": 1,
  "name": "Alexander",
  "email": "alexander@email.com",
  "created_at": "2026-04-26T09:22:31Z",
  "updated_at": "2026-04-26T09:22:31Z"
}
```

---

### Listar usuários

```bash
curl "http://localhost:8080/api/v1/users?page=1&limit=10"
```

**Resposta `200`:**
```json
{
  "data": [...],
  "total": 42,
  "page": 1,
  "limit": 10
}
```

---

### Atualizar usuário (parcial)

```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alexander Atualizado"
  }'
```

> Apenas os campos enviados são alterados. Campos omitidos permanecem inalterados.

---

### Criar pedido

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "description": "Notebook Dell XPS",
    "amount": 4500.00,
    "status": "pending"
  }'
```

**Resposta `201`:**
```json
{
  "id": 1,
  "user_id": 1,
  "description": "Notebook Dell XPS",
  "amount": 4500.00,
  "status": "pending",
  "created_at": "2026-04-26T09:22:31Z",
  "updated_at": "2026-04-26T09:22:31Z"
}
```

**Status válidos:** `pending` | `paid` | `canceled`

---

### Listar pedidos de um usuário

```bash
curl "http://localhost:8080/api/v1/users/1/orders?page=1&limit=5"
```

---

## ⚠️ Tratamento de Erros

Todos os erros seguem o formato padronizado:

```json
{
  "error": "mensagem descritiva do erro"
}
```

| Código | Situação |
|---|---|
| `400 Bad Request` | Payload inválido ou campos obrigatórios ausentes |
| `404 Not Found` | Recurso não encontrado (usuário ou pedido) |
| `409 Conflict` | E-mail já cadastrado |
| `500 Internal Server Error` | Erro inesperado no servidor |

---

## 🧠 Decisões de Arquitetura

### Por que interfaces em Repository e Service?

Cada camada depende de uma interface, nunca da implementação concreta. Isso permite:

- **Testabilidade**: criar mocks sem precisar de banco real nos testes unitários
- **Flexibilidade**: trocar PostgreSQL por outro banco alterando apenas o repository

### Por que ponteiros nos DTOs de atualização?

```go
type UpdateUserRequest struct {
    Name  *string `json:"name"`
    Email *string `json:"email"`
}
```

`nil` = campo não enviado. Sem ponteiros, seria impossível distinguir "não enviado" de "enviado como string vazia", o que causaria sobrescrita indevida de dados.

### Por que Zap em vez do `log` padrão?

O `log` padrão do Go produz texto simples. O Zap produz **JSON estruturado**, indexável em ferramentas de observabilidade como Datadog, Grafana Loki e ELK Stack — padrão em ambientes de produção.

### AutoMigrate vs Migrations versionadas

O `AutoMigrate` do GORM é conveniente para desenvolvimento — cria e atualiza tabelas automaticamente. Em produção, o recomendado é usar **golang-migrate** ou **goose** com arquivos SQL versionados (equivalente ao Flyway no ecossistema Spring Boot), garantindo controle total sobre alterações de schema e rollbacks.

---

## 👤 Autor

**Alexander** — Estudante de Engenharia de Software na UCSAL  
Especialização em backend Java/Spring Boot e Go  
[GitHub](https://github.com/alexanderbs3)