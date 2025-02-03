# 🏗️ Lari faz Croche! API-GOLANG

Este projeto é uma API desenvolvida e refatorada de uma em Java Spring, em **Golang** utilizando **Gin** como framework web, com suporte para autenticação via **JWT**, configuração com **dotenv**, e persistência de dados com **GORM** e **PostgreSQL**.

Esta API é um sistema completo de gerenciamento de produtos e categorias com autenticação de usuário(admin), com foco em **segurança** e **controle de acesso**. Ela oferece uma variedade de funcionalidades, incluindo o registro, login e a autenticação, bem como o gerenciamento de produtos em um catálogo. A API utiliza autenticação baseada em tokens JWT para garantir a segurança e o acesso restrito às rotas protegidas.

## 🚀 Tecnologias Utilizadas

- **[Golang](https://golang.org/)** – Linguagem principal do projeto
- **[Gin](https://github.com/gin-gonic/gin)** – Framework web leve e rápido
- **[GORM](https://gorm.io/)** – ORM para manipulação do banco de dados
- **[PostgreSQL Driver](https://github.com/lib/pq)** – Conexão com o banco de dados PostgreSQL
- **[JWT](https://github.com/golang-jwt/jwt)** – Autenticação segura
- **[dotenv](https://github.com/joho/godotenv)** – Gerenciamento de variáveis de ambiente
- **[Air](https://github.com/cosmtrek/air)** – Live reload para desenvolvimento

---

## 📦 Instalação

1. Clone o repositório:
   ```sh
   git clone https://github.com/seu-usuario/seu-repositorio.git
   cd seu-repositorio
   ```

2. Instale as dependências:
   ```sh
   go mod tidy
   ```

3. Configure as variáveis de ambiente:
   Crie um arquivo `.env` na raiz do projeto e adicione:
   ```ini
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=seu_usuario
   DB_PASSWORD=sua_senha
   DB_NAME=seu_banco
   JWT_SECRET=sua_chave_secreta
   ```

4. Execute as migrações do banco de dados:
   ```sh
   go run main.go migrate
   ```

---

## 🏃‍♂️ Executando a API

### Com Air (modo desenvolvimento com live reload)

```sh
air
```

### Com Go diretamente

```sh
go run main.go
```

---

## ⚙️ Estrutura do Projeto

```
├── main.go            # Ponto de entrada da aplicação
├── config             # Configurações da API (banco de dados, ambiente)
├── controllers        # Controladores da API
├── models             # Modelos de banco de dados
├── routes             # Definição de rotas
├── middleware         # Middlewares como autenticação JWT
├── utils              # Funções auxiliares
├── .env.example       # Exemplo de variáveis de ambiente
├── go.mod             # Dependências do projeto
└── README.md          # Documentação
```

---

## 🔑 Autenticação JWT

A API utiliza JWT para autenticação. Para acessar rotas protegidas, inclua o token no cabeçalho:

```
Authorization: Bearer SEU_TOKEN_AQUI
```

Para gerar um token, faça login na API enviando um `POST` para `/login` com credenciais válidas.

---

## 🛠️ Makefile para Facilitar o Desenvolvimento

Crie um arquivo `Makefile` na raiz do projeto para automatizar tarefas comuns:

```make
.PHONY: run test migrate fmt build

run:
	air

migrate:
	go run main.go migrate

test:
	go test ./...

fmt:
	go fmt ./...

build:
	go build -o app .
```

Agora você pode rodar comandos como:
- `make run` → Executa a API com live reload
- `make migrate` → Executa as migrações do banco
- `make test` → Roda os testes
- `make fmt` → Formata o código
- `make build` → Compila o binário da aplicação

---

## 📜 Licença

Este projeto está sob a licença **MIT**. Sinta-se livre para contribuir e utilizar como desejar.

---

## ✨ Contribuindo

Sinta-se à vontade para abrir **issues** e **pull requests**. Toda ajuda é bem-vinda! 🚀

