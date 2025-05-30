![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Tier](https://img.shields.io/badge/Tier-1-red)
![Postgres](https://img.shields.io/badge/Postgres-14-blue?logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-ready-blue?logo=docker)

# Ualá Twitter Backend

A scalable, clean-architecture backend for a Twitter-like system, built in Go.  
Supports in-memory and Postgres persistence.  
Includes modular domain, service, and handler layers, with robust testing and API standards.

---

## 🚀 Features

- Clean Architecture: Domain-driven, decoupled services and repositories.
- Multiple Persistence Backends: In-memory and Postgres.
- Fully-tested handlers, services, and domain logic.
- Docker-ready for local or CI/CD use.
- Endpoints for users, tweets, likes, follows, timelines, and health checks.
- Read-Optimized: Timeline endpoint fetches tweets of all followees in a single read, sorted, paginated in-memory.
- All in-memory repositories use mutexes for concurrent access.
- Decoupled Services: Handlers and use cases depend only on interfaces.
- Real-time-like timeline fetch with concurrent aggregation for scalable performance.

---

## 🛠️ Dependencies

- **Go 1.21+**
- **Postgres 14+** (Docker container included)
- **Docker / Docker Compose** (for local setup)

---

## ⚡ Getting Started

### 1. Clone and configure

```bash
git clone https://github.com/MauriGutierrez/twitter-lite.git
cd twitter-lite
go mod tidy
```

### 2. Configure Environment

Environment variables configuration.

```bash
export APP_ENV=local
export POSTGRES_DSN=postgres://postgres:postgres@localhost:5432/uala
export PORT=8080
```

### 3. Run service (Locally)

```bash
go run cmd/api/main.go
```

### 4. Run tests
```bash
go test ./...
```

### 5. API Examples (using localhost:8080)

### Create User

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Mauricio","document":"38307207"}'
```

```bash
Sample response:
{"id": "usr_38307207"}
```

### Post Tweet

```bash

curl -X POST http://localhost:8080/tweets \
  -H "X-User-ID: usr_38207274" \
  -H "Content-Type: application/json" \
  -d '{"content":"Soy Mauri y este es mi primer tweet?"}'
```

```bash
Sample response:
{"id": "f7bda8a9-1234-4567-890a-b2e1df789abc"}
```

### Follow User

```bash
curl -X POST http://localhost:8080/follow \
  -H "X-User-ID: usr_38207209" \
  -H "Content-Type: application/json" \
  -d '{"followee_id":"usr_38207274"}'
```

```bash
Sample response:
{}
```
### Get Timeline

```bash
curl -X GET http://localhost:8080/timeline \
  -H "X-User-ID: usr_38207209"
```

```bash
Sample response:
[
  {
    "id": "a1b2c3d4-e5f6-7890-1234-5678abcdef90",
    "user_id": "usr_38207274",
    "content": "Soy Mauri y este es mi primer tweet?",
    "likes": 1,
    "created_at": "2025-05-29T18:23:12-03:00"
  }
]
```

### Like Tweet

```bash
curl -X POST http://localhost:8080/tweets/{tweet_id}/like \
  -H "X-User-ID: usr_38207274"
```  

```bash
Sample response:
Response: 204 No Content
```

### Health Check

```bash
curl -X GET http://localhost:8080/health
```
    
```bash
Sample response:
{
  "env": "local",
  "name": "uala-twitter",
  "version": "1.0.0"
}
```

## 6. Production-Ready Considerations

- **Users:** Postgres or other relational DB for transactional integrity and uniqueness constraints.

- **Tweets:** Sharded relational DB, or scalable NoSQL like DynamoDB for horizontal scaling and high write throughput.

- **Follows/Likes:** Redis for efficient set operations (who follows/liked).

- **Read Optimization:** Redis/memcached for timeline caching; denormalized timeline tables for high-velocity reads.

- **Scaling:** Message queues (Kafka/SQS) for async fan-out, background jobs for heavy write operations.

