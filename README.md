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
git clone https://github.com/MauriGutierrez/uala-twitter-backend.git
cd uala-twitter-backend
go mod tidy
