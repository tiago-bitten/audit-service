# Audit Service

Serviço de registro de logs de auditoria multi-projeto.

## Rodando

```bash
docker-compose up
```

Variáveis de ambiente:

| Variável | Padrão | Descrição |
|---|---|---|
| `MONGO_URL` | — | Connection string do MongoDB (obrigatória) |
| `MONGO_DATABASE` | `auditservice` | Nome do banco |
| `PORT` | `8080` | Porta HTTP |

## API

### 1. Criar projeto

```
POST /v1/projects
```

```json
{ "name": "meu-projeto" }
```

Resposta `201`:

```json
{ "success": true, "project_id": "b1c2d3e4-..." }
```

Guarde o `project_id` — ele será usado como header de autenticação nas demais rotas.

### 2. Registrar log

```
POST /v1/auditlogs
X-Project-Id: b1c2d3e4-...
```

```json
{
  "message": "usuário atualizou perfil",
  "date": "2026-02-15T10:30:00Z",
  "item_id": "user-123",
  "user_id": "admin-1",
  "group_id": "org-456",
  "object": { "field": "email", "old": "a@x.com", "new": "b@x.com" }
}
```

Campos obrigatórios: `message`, `date` (UTC).
Campos opcionais: `item_id`, `user_id`, `group_id`, `object`.

Resposta `201`:

```json
{ "success": true }
```

### 3. Consultar logs

```
GET /v1/auditlogs
X-Project-Id: b1c2d3e4-...
```

Query params (todos opcionais):

| Param | Descrição |
|---|---|
| `item_id` | Filtrar por item |
| `user_id` | Filtrar por usuário |
| `group_id` | Filtrar por grupo |
| `limit` | Quantidade por página (padrão: 50) |
| `offset` | Pular N registros |

Resposta `200`:

```json
{
  "success": true,
  "data": [
    {
      "id": "67ab...",
      "message": "usuário atualizou perfil",
      "date": "2026-02-15T10:30:00Z",
      "project_id": "b1c2d3e4-...",
      "item_id": "user-123",
      "user_id": "admin-1",
      "group_id": "org-456",
      "object": { "field": "email", "old": "a@x.com", "new": "b@x.com" }
    }
  ]
}
```
