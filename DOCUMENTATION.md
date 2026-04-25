# Project Architecture & Documentation

This document serves as a deep dive into the technical architecture, design patterns, and features implemented in the Multi-Tenant SaaS platform. 

## 1. Architectural Pattern: Layered (Clean) Architecture

The backend is structured using a Layered Architecture. This ensures that concerns are separated, making the code testable, maintainable, and scalable.

### Directory Breakdown
- **`cmd/main.go`**: The entry point. It wires up the database, initializes the repositories and services, runs database migrations, starts the background worker pool, configures the router, and starts the HTTP server.
- **`internal/models/`**: Defines the Core Domain (Data Structures). Examples: `User`, `Organization`, `Project`, `Task`, and `Membership`.
- **`internal/repository/`**: The Data Layer. This is the **only** layer that interacts with PostgreSQL. It executes raw SQL queries using the `pgxpool` library for high-performance connection pooling.
- **`internal/services/`**: The Business Logic layer. It sits between the Handlers and the Repositories. Handlers pass requests here, and the Service decides what rules to apply before asking the Repository to save/fetch data.
- **`internal/handlers/`**: The Transport layer (HTTP). It receives HTTP requests, parses JSON bodies, extracts User IDs from context, calls the Service layer, and formats the HTTP JSON response.
- **`internal/middleware/`**: Contains HTTP interceptors. Handles CORS, JWT authentication, and RBAC (Role-Based Access Control).
- **`internal/worker/`**: Contains the Go-native asynchronous background job processing engine.
- **`internal/utils/`**: Helper functions (e.g., standardizing JSON responses, validating emails/passwords, hashing passwords, generating JWTs).

---

## 2. Core Features & Implementation Details

### A. Deep Multi-Tenancy (Data Isolation)
In a SaaS application, users belong to "Organizations" (Tenants). It is a catastrophic security failure if User A can see Organization B's data.
**How it is solved:**
- We do not just trust the user ID in the handler. Every single database query in the `repository` layer (e.g., `GetOrganizations`, `GetProjects`) performs a strict SQL `JOIN` with the `memberships` table. 
- Example: When fetching projects, the query enforces `WHERE memberships.user_id = $1`. If the user is not explicitly a member of the organization that owns the project, the database returns absolutely nothing.

### B. Authentication & Security (JWT)
- **Login Flow**: When a user logs in, `internal/utils/jwt.go` generates a JSON Web Token (JWT) signed with a secret key. This token contains the user's `ID`. We implemented both Access Tokens (short-lived) and Refresh Tokens (long-lived) for optimal security.
- **Middleware**: Every protected route goes through `middleware.Auth()`. This intercepts the request, validates the `Bearer` token in the `Authorization` header, extracts the `UserID`, and injects it into the HTTP Request Context (`r.Context()`) so handlers can use it securely.

### C. Role-Based Access Control (RBAC)
- **Concept**: Even within an Organization, some users are `Admins` and some are `Members`.
- **Implementation**: The `middleware.RequireRole("admin")` function intercepts requests that alter organizational state (like inviting members). It queries the database to check the user's specific role in that specific organization before allowing the HTTP request to reach the handler.

### D. Asynchronous Background Workers
- **The Problem**: Actions like sending emails or generating PDFs take seconds to complete. If done synchronously, the user's browser hangs waiting for the API to respond.
- **The Solution**: In `internal/worker/pool.go`, we built a Thread Pool. 
  - On startup, the server spins up 5 Go `workers` (Goroutines) that infinitely listen to a Go `channel` (queue).
  - When a task is created (`TaskHandler`), we instantly drop a `TaskCreatedJob` into the queue and return a `201 Created` HTTP response to the user in 0ms.
  - The background worker picks up the job from the queue and processes it asynchronously.

### E. Database Migrations
- **The Problem**: Manually running SQL scripts to alter tables causes team chaos and breaking changes in production.
- **The Solution**: We integrated `golang-migrate/migrate`. The `db/migrations/` folder contains versioned `.up.sql` and `.down.sql` files. 
- In `cmd/main.go`, before the server even opens port 8081, it compares the migration files against the `schema_migrations` table in PostgreSQL. If there are new SQL files, it executes them automatically.

### F. Automated DevOps (Docker & CI/CD)
- **Testing**: We use Table-Driven testing (`internal/utils/validator_test.go`) to test multiple edge cases without duplicating code.
- **GitHub Actions (`.github/workflows/ci.yml`)**: On every push, GitHub spins up a Linux server, installs Go, pulls the code, and runs the tests to ensure nothing was broken.
- **Docker (`Dockerfile` & `docker-compose.yml`)**: The backend is packaged into an isolated, lightweight Alpine container. Anyone can type `docker-compose up` to instantly provision a perfectly configured PostgreSQL database and the Go API without installing dependencies.

### G. Swagger API Documentation
- By adding declarative comments (`// @Summary`, `// @Param`) above our handler functions and running `swag init`, we generate a `docs/` module.
- `http-swagger` serves this JSON as an interactive webpage at `/swagger/index.html`, allowing frontend engineers to see exactly what endpoints exist and what JSON bodies they require.

---

## 3. The Frontend (React)
The frontend (`saas-frontend`) is built with modern React.
- **State Management**: Uses React Context/Hooks to manage user sessions.
- **Styling**: Vanilla CSS utilizing modern CSS variables for a premium, consistent design system (Glassmorphism, gradients, hover micro-animations).
- **Integration**: The frontend handles JWT tokens gracefully, attaching them to the `Authorization` header for every backend fetch request, and uses Optional Chaining (`?.`) to safely handle asynchronous data states without crashing.
