# 🚀 Deploy no Render - Lariharumi Croche Backend

## 📋 Configuração para Render

### 🔧 Variáveis de Ambiente Necessárias

Configure as seguintes variáveis de ambiente no painel do Render:

```bash
# Configurações do Banco de Dados
DB_HOST=seu-host-postgresql
DB_USER=seu-usuario-postgresql
DB_PASSWORD=sua-senha-postgresql
DB_NAME=nome-do-banco
DB_PORT=5432

# Configurações de Segurança
JWT_SECRET=seu-jwt-secret-muito-seguro-aqui

# Configurações da API ImgBB
IMGBB_API_KEY=sua-chave-api-imgbb

# Configurações do Frontend
FRONTEND_URL=https://seu-frontend-url.com
BASEURL=https://sua-app-backend.onrender.com

# Configurações do Redis (opcional)
REDIS_URL=sua-url-redis
REDIS_PASSWORD=sua-senha-redis

# Porta (Render define automaticamente)
PORT=10000
```

### 🏗️ Configuração do Build

O Render usará automaticamente:
- **Build Command**: `go build -o app .`
- **Start Command**: `./app`

### 🔄 Cronjob Automático

A aplicação inclui um cronjob interno que:
- ✅ Faz ping a cada 25 segundos para manter a aplicação ativa
- ✅ Usa endpoints públicos (`/health`, `/ping`, `/categories`, `/products`)
- ✅ Evita que o Render derrube a aplicação por inatividade
- ✅ Funciona automaticamente sem configuração adicional

### 📊 Endpoints de Monitoramento

- **`GET /health`** - Health check completo da aplicação
- **`GET /ping`** - Ping simples para manter ativo
- **`GET /categories`** - Endpoint público usado pelo cronjob
- **`GET /products`** - Endpoint público usado pelo cronjob

### 🚨 Troubleshooting

#### Problema: Aplicação sendo derrubada após 30 segundos
**Solução**: O cronjob interno já resolve isso automaticamente. Verifique os logs para confirmar que está funcionando:
```
🔄 Cronjob de health check iniciado - ping a cada 25 segundos
✅ Ping bem-sucedido para /categories (status: 200)
```

#### Problema: Erro de conexão com banco
**Solução**: Verifique se as variáveis `DB_*` estão configuradas corretamente.

#### Problema: Erro de Redis
**Solução**: Redis é opcional. A aplicação funciona sem ele (graceful degradation).

### 📝 Logs Importantes

Monitore estes logs no Render:
- `Servidor rodando na porta 10000....`
- `🔄 Cronjob de health check ativo`
- `✅ Ping bem-sucedido para /endpoint`
- `Banco de dados conectado!`
- `Redis conectado com sucesso!` (se Redis estiver configurado)

### 🔗 URLs Importantes

Após o deploy, sua aplicação estará disponível em:
- **API Base**: `https://sua-app-backend.onrender.com`
- **Health Check**: `https://sua-app-backend.onrender.com/health`
- **Ping**: `https://sua-app-backend.onrender.com/ping`

### ⚡ Performance

Com as melhorias implementadas:
- **Cache Redis**: Respostas mais rápidas
- **Soft Delete**: Dados seguros
- **Upload Assíncrono**: Uploads paralelos
- **Paginação Otimizada**: Menos dados transferidos
- **Cronjob**: Aplicação sempre ativa

### 🎯 Próximos Passos

1. Configure as variáveis de ambiente
2. Faça o deploy
3. Teste os endpoints públicos
4. Monitore os logs do cronjob
5. Configure seu frontend para usar a nova URL

---

**✨ Sua aplicação estará sempre ativa no Render com o cronjob interno!**
