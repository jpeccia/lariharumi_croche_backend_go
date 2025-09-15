# 🏗️ Lari faz Croche! API-GOLANG

Este projeto é uma API desenvolvida e refatorada de uma em Java Spring, em **Golang** utilizando **Gin** como framework web, com suporte para autenticação via **JWT**, configuração com **dotenv**, e persistência de dados com **GORM** e **PostgreSQL**.

Esta API é um sistema completo de gerenciamento de produtos e categorias com autenticação de usuário(admin), com foco em **segurança** e **controle de acesso**. Ela oferece uma variedade de funcionalidades, incluindo o registro, login e a autenticação, bem como o gerenciamento de produtos em um catálogo. A API utiliza autenticação baseada em tokens JWT para garantir a segurança e o acesso restrito às rotas protegidas.

## 🚀 Tecnologias Utilizadas

- **[Golang](https://golang.org/)** – Linguagem principal do projeto
- **[Gin](https://github.com/gin-gonic/gin)** – Framework web leve e rápido
- **[GORM](https://gorm.io/)** – ORM para manipulação do banco de dados
- **[PostgreSQL Driver](https://github.com/lib/pq)** – Conexão com o banco de dados PostgreSQL
- **[Redis](https://redis.io/)** – Sistema de cache para melhor performance
- **[JWT](https://github.com/golang-jwt/jwt)** – Autenticação segura
- **[dotenv](https://github.com/joho/godotenv)** – Gerenciamento de variáveis de ambiente
- **[Air](https://github.com/cosmtrek/air)** – Live reload para desenvolvimento
- **[ImgBB API](https://imgbb.com/)** – Upload de imagens para nuvem

---

## 📦 Instalação

1. Clone o repositório:
   ```sh
   git clone https://github.com/jpeccia/lariharumi_croche_backend_go
   cd lariharumi_croche_backend_go
   ```

2. Instale as dependências:
   ```sh
   go mod tidy
   ```

3. Configure as variáveis de ambiente:
   Crie um arquivo `.env` na raiz do projeto e adicione:
   ```ini
   # Configurações do Banco de Dados
   DB_HOST=localhost
   DB_USER=postgres
   DB_PASSWORD=password
   DB_NAME=lariharumi_croche
   DB_PORT=5432

   # Configurações do Redis
   REDIS_URL=localhost:6379
   REDIS_PASSWORD=

   # Configurações de Segurança
   JWT_SECRET=seu_jwt_secret_aqui_muito_seguro

   # Configurações da API ImgBB
   IMGBB_API_KEY=sua_chave_api_imgbb_aqui

   # Configurações do Frontend
   FRONTEND_URL=http://localhost:3000
   BASEURL=http://localhost:8080
   ```

4. Rode seu docker para criar o banco de dados:
   ```sh
   docker-compose up -d
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

## ✨ Funcionalidades Implementadas

### 🗃️ **Soft Delete**
- Todos os modelos (User, Category, Product) agora usam soft delete
- Registros não são removidos permanentemente do banco
- Funções `HardDelete*` disponíveis para remoção permanente (apenas admin)

### ⚡ **Sistema de Cache com Redis**
- Cache automático para produtos e categorias
- TTL configurável (15min para produtos, 1h para categorias)
- Invalidação automática quando dados são modificados
- Funciona mesmo sem Redis (graceful degradation)

### 🔍 **Paginação Melhorada**
- Metadados completos de paginação (total, páginas, navegação)
- Contagem total de registros
- Limites de segurança (máximo 100 por página)
- Suporte a pesquisa com paginação

### 🚀 **Upload Assíncrono**
- Upload de múltiplas imagens em paralelo
- Pool de workers configurável
- Monitoramento de progresso
- Tratamento de erros individual por arquivo

### 🎯 **Otimizações de Performance**
- Preload automático de relacionamentos (evita N+1 queries)
- Cache inteligente com invalidação seletiva
- Upload paralelo de imagens
- Queries otimizadas com índices

### 🔄 **Cronjob para Render**
- Ping automático a cada 25 segundos para manter aplicação ativa
- Evita que o Render derrube a aplicação por inatividade
- Endpoints de health check (`/health`, `/ping`)
- Funciona automaticamente sem configuração adicional

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
.PHONY: run test docker fmt build

run:
	air

docker:
	docker-compose up -d

test:
	go test ./...

fmt:
	go fmt ./...

build:
	go build -o app .
```

Agora você pode rodar comandos como:
- `make run` → Executa a API com live reload
- `make docker` → Executa as migrações do banco
- `make test` → Roda os testes
- `make fmt` → Formata o código
- `make build` → Compila o binário da aplicação

---

## 📜 Licença

Este projeto está sob a licença **MIT**. Sinta-se livre para contribuir e utilizar como desejar.

---

## 🚀 Deploy no Render

A aplicação está preparada para deploy no Render com cronjob automático para manter a aplicação ativa.

### 📋 Configuração Rápida

1. **Conecte seu repositório** no Render
2. **Configure as variáveis de ambiente** (veja `RENDER_DEPLOY.md`)
3. **Deploy automático** - o cronjob manterá a aplicação ativa

### 🔄 Cronjob Automático

- ✅ Ping a cada 25 segundos (menos que 30s do Render)
- ✅ Usa endpoints públicos (`/health`, `/ping`, `/categories`)
- ✅ Funciona automaticamente sem configuração
- ✅ Evita que o Render derrube a aplicação

### 📊 Endpoints de Monitoramento

- `GET /health` - Health check completo
- `GET /ping` - Ping simples
- `GET /categories` - Usado pelo cronjob
- `GET /products` - Usado pelo cronjob

**📖 Para instruções detalhadas, veja `RENDER_DEPLOY.md`**

---

## ✨ Contribuindo

Sinta-se à vontade para abrir **issues** e **pull requests**. Toda ajuda é bem-vinda! 🚀

