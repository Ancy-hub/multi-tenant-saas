# Multi-Tenant SaaS Backend 🚀

A production-ready, enterprise-grade backend for a Multi-Tenant SaaS application. Built with **Go** and **PostgreSQL**, this project demonstrates modern scalable architecture, deep data isolation, advanced concurrency, and robust Dev/Ops practices.

---

## ✨ Key Features

- **Full-Stack Implementation**: Features a sleek, modern React frontend built for a premium user experience, seamlessly integrated with the Go backend.
- **Strict Multi-Tenancy**: Deep data isolation at the repository layer. Users can only access organizations, projects, and tasks they are explicitly members of.
- **Advanced Security**: Secure authentication flow using JSON Web Tokens (JWT) with Access and Refresh tokens for session management.
- **Role-Based Access Control (RBAC)**: Custom middleware enforcing granular permissions (`Admin`, `Member`) across all endpoints.
- **Asynchronous Background Workers**: Implements a highly concurrent Go worker pool using native Channels and Goroutines to process non-blocking background tasks (e.g., email dispatch simulations) without hanging the API.
- **Automated Database Migrations**: Schema versioning managed seamlessly through `golang-migrate`, embedded directly into the Go binary to automatically sync the database state on startup.
- **Docker Containerization**: Includes a multi-stage Docker build resulting in a secure, lightweight Alpine image, orchestrated locally via `docker-compose`.
- **Automated CI/CD Pipeline**: GitHub Actions workflow that automatically lints, builds, and runs table-driven unit tests on every push.
- **Swagger API Documentation**: Interactive, auto-generated OpenAPI documentation accessible via the browser.

## 🛠️ Technology Stack

- **Backend**: Go (1.22+), go-chi/chi
- **Frontend**: React.js, HTML5, Vanilla CSS
- **Database**: PostgreSQL 15, `pgxpool` (jackc/pgx/v5)
- **Security**: JWT (Access & Refresh Tokens)
- **DevOps**: Docker, Docker Compose, GitHub Actions
- **Migrations**: `golang-migrate/migrate/v4`
- **Documentation**: `swaggo/swag`

---

## 🚀 Getting Started

There are two ways to run this application: **Docker** (Recommended) or **Local Development**.

### Option 1: Run with Docker (The easy way)
Ensure you have Docker and Docker Compose installed.

```bash
# Spin up the PostgreSQL database and the Go Backend
docker-compose up --build
```
*The API will be available at `http://localhost:8081`.*

### Option 2: Run locally for Development
You will need Go 1.22+, Node.js, and a local PostgreSQL instance running.

1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/multi-tenant-saas.git
   cd multi-tenant-saas
   ```

2. **Run the Backend:**
   ```bash
   export DATABASE_URL="postgres://postgres:root@localhost:8080/postgres?sslmode=disable"
   make run
   # Or manually: go run cmd/main.go
   ```

3. **Run the Frontend:**
   Open a new terminal window:
   ```bash
   cd saas-frontend
   npm install
   npm start
   ```

---

## 📖 API Documentation (Swagger)

Once the server is running, you can explore and interact with the endpoints through the auto-generated Swagger UI.

👉 **Visit:** `http://localhost:8081/swagger/index.html`

## 🧪 Testing

The project uses Go's idiomatic "Table-Driven Testing" approach.

```bash
# Run all unit tests
make test
```

## 🏗️ Architecture Overview

- **`cmd/`**: Entry point of the application (`main.go`). Handles configuration, DB connection, routing, and migration execution.
- **`internal/handlers/`**: HTTP transport layer. Parses JSON, extracts JWT claims, and formats responses.
- **`internal/services/`**: Business logic layer.
- **`internal/repository/`**: Database interaction layer. All multi-tenant data isolation `JOIN` queries happen here using `pgxpool`.
- **`internal/middleware/`**: JWT validation and RBAaC enforcement.
- **`internal/worker/`**: The asynchronous Worker Pool logic and Job definitions.
- **`db/migrations/`**: Raw SQL up/down files for schema evolution.

---

*Designed and engineered by an Automation Engineer transitioning into Backend Software Engineering.*
