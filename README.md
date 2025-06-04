# 🚀 Go API com Swagger - Projeto Base (MongoDB)

Este é um projeto base de uma API RESTful desenvolvida em **Go (Golang)**, com arquitetura em camadas (`controller`, `service`, `repository`) e documentação automática com **Swagger** utilizando o `swaggo`.

Banco de dados utilizado: **MongoDB** (com persistência via Docker).

---

## ✨ Funcionalidades

- ✅ Cadastro e listagem de usuários
- ✅ Autenticação JWT com login
- ✅ Validação de campos obrigatórios
- ✅ Verificação de e-mail duplicado (índice único no Mongo)
- ✅ Documentação Swagger (`/doc/api`)
- ✅ Testes automatizados com `net/http/httptest`
- ✅ Arquitetura escalável com separação de responsabilidades
- ✅ Pronto para rodar com Docker

---

## 🧱 Estrutura do Projeto

```
go-api/
├── config/         # Carregamento de variáveis de ambiente
├── controllers/    # Lida com requisições HTTP
├── database/       # Conexão com MongoDB e criação de índices
├── docs/           # Arquivos gerados pelo swag
├── models/         # Estruturas dos dados (User)
├── repositories/   # Interação com o banco de dados
├── routes/         # Definição de rotas
├── services/       # Regras de negócio (validações, etc)
├── tests/          # Testes automatizados
├── main.go         # Ponto de entrada da aplicação
├── .env            # Variáveis de ambiente
└── docker-compose.yml
```

---

## 📚 Documentação Swagger

Geração automática com `swaggo/swag`.

> Acesse após subir a aplicação:

```
http://localhost:8080/doc/api
```

---

## 🛠️ Como rodar o projeto

### 1. Clone o repositório

```bash
git clone https://github.com/seu-usuario/go-api.git
cd go-api
```

### 2. Crie o arquivo `.env`

```env
MONGO_URI=mongodb://mongo:27017
MONGO_DB=goapidb
PORT=8080
```

> Se estiver rodando local sem Docker, use:
> `MONGO_URI=mongodb://localhost:27017`

### 3. Gere a documentação Swagger

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

### 4. Suba os containers (MongoDB + API)

```bash
docker-compose up --build
```


### 5. Rodando sem Docker (opcional)

Caso prefira executar a API diretamente com o Go instalado na sua máquina,
certifique-se de que um servidor MongoDB está ativo em `localhost:27017` e que o
arquivo `.env` possui a variável `MONGO_URI` apontando para essa instância.

```bash
go mod download
go run main.go
```


### 6. Cadastro rápido via CLI

É possível criar usuários diretamente pelo terminal para facilitar o desenvolvimento:

```bash
go run ./cmd/createuser -name "Admin" -email admin@example.com -password 123456 -admin
```

O parâmetro `-admin` é opcional e cria um usuário administrador.


---

## ✅ Exemplos de Endpoints

### `GET /users`

Retorna todos os usuários cadastrados.

### `POST /users`

Cria um novo usuário:

```json
{
  "name": "Josuel",
  "email": "josuel@example.com",
  "password": "suaSenha"
}
```

### `POST /login`

Retorna um token JWT:

```json
{
  "email": "josuel@example.com",
  "password": "suaSenha"
}
```

---

## 🧪 Rodando os Testes

```bash
go test ./tests -v
```

> Os testes verificam:
> - Criação com sucesso
> - Erro ao criar com e-mail duplicado
> - Erro ao criar com campos vazios

---

## 📦 Tecnologias Utilizadas

- **Go 1.22+**
- **MongoDB** (com driver oficial)
- **Gorilla Mux**
- **swaggo/swag** (documentação)
- **Docker** + **Docker Compose**
- **httptest** (testes de integração)

---

## 🧠 Créditos

Desenvolvido por [@razordz](https://github.com/razordz) com 💻 e ☕.
