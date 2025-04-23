
# 📦 Go API com Swagger - Projeto Base

Este projeto é uma API RESTful simples desenvolvida em **Go (Golang)**, utilizando o padrão de camadas (`controller`, `service`, `repository`) e documentada com **Swagger** via `swaggo`.

---

## 🚀 Funcionalidades

- ✅ Cadastro de usuários
- ✅ Listagem de usuários
- ✅ Validação de campos
- ✅ Integração com banco PostgreSQL via GORM
- ✅ Documentação automática com Swagger (`/doc/api`)
- ✅ Arquitetura escalável com camadas organizadas

---

## 🧱 Estrutura do Projeto

```
go-api/
├── config/         # Carrega variáveis de ambiente (.env)
├── controllers/    # Controladores HTTP
├── database/       # Inicialização e conexão com o banco
├── models/         # Estrutura dos dados (User)
├── repositories/   # Acesso ao banco
├── routes/         # Rotas e endpoints
├── services/       # Regras de negócio
├── .env            # Configurações sensíveis
├── main.go         # Ponto de entrada
└── README.md       # Este arquivo
```

---

## 📚 Documentação Swagger

Após rodar a aplicação, acesse:

👉 [`http://localhost:8080/doc/api`](http://localhost:8080/doc/api)

---

## 🛠️ Como rodar o projeto

### 1. Clone o repositório

```bash
git clone https://github.com/seu-usuario/go-api.git
cd go-api
```

### 2. Instale as dependências

```bash
go mod tidy
```

### 3. Configure o `.env`

Crie um arquivo `.env` com base no modelo:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=123456
DB_NAME=goapidb
PORT=8080
```

### 4. Gere a documentação Swagger

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

### 5. Rode a aplicação

```bash
go run main.go
```

---

## ✅ Exemplos de endpoints

### `GET /users`

Retorna todos os usuários cadastrados.

### `POST /users`

Cria um novo usuário:

```json
{
  "name": "Josuel",
  "email": "josuel@example.com"
}
```

---

## 🧪 Testes

Execute os testes com:

```bash
go test ./tests -v
```

---

## 📦 Tecnologias utilizadas

- Go 1.22+
- GORM
- Gorilla Mux
- swaggo/swag
- PostgreSQL
- Docker (opcional)
