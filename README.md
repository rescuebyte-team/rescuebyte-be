# RescueByte Backend

![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)

## Overview

RescueByte Backend is a high-performance backend application built with a focus on maintainability, scalability, and robust application security. RescueByte powers a surplus food and goods rescue marketplace connecting merchants and stores with consumers through mystery boxes and discounted rescue pricing.

The project structure adheres to modern software engineering practices, ensuring that the codebase remains clean, modular, scalable, and consistent.

Core technologies used:

* [Go v1.26.0](https://go.dev/doc/go1.26)
* [PostgreSQL](https://www.postgresql.org/)
* [Redis](https://redis.io/)
* [Fiber v3](https://docs.gofiber.io/) as the HTTP web framework
* [pgx/v5](https://github.com/jackc/pgx) as the PostgreSQL driver
* [sqlx](https://github.com/jmoiron/sqlx) for database operations and mapping
* [Squirrel](https://github.com/Masterminds/squirrel) as the fluent SQL query builder

---

# Project Setup

## Running Locally

> **Note:** This project uses [Task](https://taskfile.dev/) as its task runner. All commands can also be executed manually by checking the contents of `Taskfile.yml`.

### 1. Ensure prerequisite dependencies are installed

Install Go (version 1.26 or higher) and Task on your machine:

```bash
go version && task --version # windows
go version && task --version # unix
```

### 2. Set up environment variables

Copy the example environment configuration:

```bash
cp ./config/.env.example ./config/.env
```

Adjust the configuration values in `./config/.env` as needed (PostgreSQL connection, Redis, JWT Secret, Midtrans, Supabase, Google OAuth).

### 3. Install all dependencies and development tools

This command runs `go mod tidy` and installs required CLI tools (`air`, `golang-migrate`, and `swag`):

```bash
task
# or
task install
```

### 4. Configure Air (Hot-Reload)

Air configuration is pre-configured in `.air.toml`. Adjust the target binary build command based on your OS if necessary:

```toml
[build]
cmd = "go build -o ./tmp/main ./cmd/app"     # unix / linux / mac
cmd = "go build -o ./tmp/main.exe ./cmd/app" # windows
```

### 5. Run database migrations

Apply all database migrations to PostgreSQL:

```bash
task migrate:up
```

### 6. Generate Swagger documentation (Optional)

When adding or updating endpoints or DTO structs, regenerate the Swagger API docs:

```bash
task swagger
```

### 7. Run the application in development mode

```bash
task dev
```

The HTTP server will run by default at `http://localhost:8080`.

---

## Running with Docker

### 1. Install Docker and Docker Compose

Ensure Docker and Docker Compose are installed and running:

```bash
docker --version && docker compose version
```

### 2. Prepare the `.env` file

Copy the environment template:

```bash
cp ./config/.env.example ./config/.env
```

Configure your environment variables. When running within Docker network, ensure:
* `DB_HOST=db`
* `REDIS_ADDR=redis:6379`

### 3. Build and start services

Start all services (`app`, `db`, `redis`) in detached mode:

```bash
task docker:up
# or rebuild containers:
task docker:build
```

### 4. Run database migrations

To execute database migrations from your host machine against the Docker container:

```bash
task migrate:up
```

Make sure `DB_HOST=localhost` in your local environment file when executing migrations from outside the Docker network.

### 5. Monitor logs and container management

To view real-time logs from Docker services:

```bash
task docker:logs
```

To stop and remove containers:

```bash
task docker:down
```

---

# Documentation

* Interactive API documentation is available via Swagger UI at [http://localhost:8080/swagger](http://localhost:8080/swagger) or [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html).
* Raw OpenAPI specifications are accessible at [JSON](http://localhost:8080/docs/swagger.json) and [YAML](http://localhost:8080/docs/swagger.yaml).
* Detailed architectural guidelines and coding standards are documented in [docs/ARCHITECTURE_GUIDE.md](docs/ARCHITECTURE_GUIDE.md).

---

# Architecture

The system follows a combination of modern architectural patterns and principles:

* Clean Architecture
* Feature-Sliced Layered Architecture
* Domain-Driven Design (DDD)
* Hexagonal Architecture (Ports & Adapters)
* SOLID Principles

---

# Key Features

* **Database Migrations** using [golang-migrate](https://github.com/golang-migrate/migrate)
* **Request Validation** powered by [validator v10](https://github.com/go-playground/validator) with custom domain validation rules
* **Structured Logging** using [zerolog](https://github.com/rs/zerolog)
* **Centralized Error Handling** with uniform HTTP response formatting
* **Authentication & RBAC** using [JWT](https://github.com/golang-jwt/jwt) and [crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
* **Google OAuth2** for seamless social login and auto-registration
* **Payment Gateway Integration** with Midtrans Snap API and webhook notification handling verified with SHA-512 signatures
* **Cloud Object Storage** integrated with [Supabase Storage](https://supabase.com/docs/guides/storage)
* **Caching & Key-Value Store** using [go-redis/v9](https://github.com/redis/go-redis)
* **Atomic Stock Management & Transactions** utilizing a database Transactor (`WithinTransaction`) and conditional SQL updates to prevent race conditions
* **Interactive API Documentation** via [Swagger / OpenAPI](https://github.com/swaggo/swag) and [gofiber/contrib/swaggo](https://github.com/gofiber/contrib)
* **Environment Configuration** managed with [Viper](https://github.com/spf13/viper)
* **Security & Utility Middlewares** including CORS, Panic Recovery, Request ID tracing, and HTTP request logging
* **Docker & Docker Compose Integration** with multi-stage builds and containerized orchestration

---

# Modules & API Endpoints

Base URL: `http://localhost:8080/api/v1`

### 1. Authentication (`/auth`)
| Method | Endpoint | Description | Access |
| :--- | :--- | :--- | :--- |
| `POST` | `/auth/register` | Register a new user account | Public |
| `POST` | `/auth/login` | Authenticate using email and password | Public |
| `POST` | `/auth/refresh` | Rotate access token using refresh token | Public |
| `POST` | `/auth/logout` | Revoke refresh token and invalidate session | Public |
| `GET` | `/auth/google` | Obtain Google OAuth2 consent URL | Public |
| `GET` | `/auth/google/callback` | Google OAuth2 callback (auto-register / login) | Public |
| `GET` | `/auth/me` | Fetch authenticated user's session profile | Protected (`RequireAuth`) |

### 2. User Profile (`/users`)
| Method | Endpoint | Description | Access |
| :--- | :--- | :--- | :--- |
| `GET` | `/users/profile` | Retrieve current user profile details | Protected (`RequireAuth`) |
| `PUT` | `/users/profile` | Update profile information (name, avatar, phone, gender, dob) | Protected (`RequireAuth`) |
| `PATCH` | `/users/role` | Update user role | Protected (Admin) |

### 3. System Enums Lookup (`/enums`)
| Method | Endpoint | Description | Access |
| :--- | :--- | :--- | :--- |
| `GET` | `/enums` | Retrieve full system enum dictionary (filterable with `?types=...`) | Public |

### 4. Cloud Storage (`/storage`)
| Method | Endpoint | Description | Access |
| :--- | :--- | :--- | :--- |
| `POST` | `/storage/upload` | Upload a single file/image to Supabase Storage | Protected (`RequireAuth`) |
| `POST` | `/storage/upload-multiple` | Upload multiple files simultaneously | Protected (`RequireAuth`) |

### 5. Merchant Stores (`/stores`)
| Method | Endpoint | Description | Access |
| :--- | :--- | :--- | :--- |
| `GET` | `/stores` | List partner stores (filter by name, city, geolocation radius) | Public |
| `GET` | `/stores/:id` | Get store details by ID | Public |
| `POST` | `/stores` | Register a new partner store | Protected (`RequireAuth`) |
| `GET` | `/stores/me/profile` | Get store profile of the authenticated merchant | Protected (`RequireAuth`) |
| `PUT` | `/stores/me/profile` | Update merchant store details | Protected (`RequireAuth`) |
| `PUT` | `/stores/me/operational` | Update store operational hours and working days | Protected (`RequireAuth`) |
| `PATCH` | `/stores/me/status` | Toggle store status (open / closed) | Protected (`RequireAuth`) |

### 6. Products & Marketplace Cards (`/products`)
| Method | Endpoint | Description | Access |
| :--- | :--- | :--- | :--- |
| `GET` | `/products` | List mystery box product cards (search, city, latitude, longitude) | Public |
| `GET` | `/products/:id` | Get complete product details with variants and photo gallery | Public |
| `POST` | `/products` | Create a new mystery box product with variants and images | Protected (`RequireAuth`) |
| `GET` | `/products/me/list` | List products owned by the merchant's store | Protected (`RequireAuth`) |
| `PUT` | `/products/:id` | Update product details and variants | Protected (`RequireAuth`) |
| `PATCH` | `/products/:id/status` | Toggle product active status | Protected (`RequireAuth`) |
| `DELETE` | `/products/:id` | Soft delete a product | Protected (`RequireAuth`) |

### 7. Shopping Cart (`/cart`)
| Method | Endpoint | Description | Access |
| :--- | :--- | :--- | :--- |
| `GET` | `/cart` | Retrieve shopping cart items grouped by store with discounts & distance | Protected (`RequireAuth`) |
| `GET` | `/cart/count` | Get total item quantity in user's cart | Protected (`RequireAuth`) |
| `POST` | `/cart/items` | Add product variant to shopping cart | Protected (`RequireAuth`) |
| `PATCH` | `/cart/items/:id` | Update item quantity in shopping cart | Protected (`RequireAuth`) |
| `DELETE` | `/cart/items/:id` | Remove a single item from shopping cart | Protected (`RequireAuth`) |
| `DELETE` | `/cart/stores/:store_id` | Remove all items from a specific store | Protected (`RequireAuth`) |
| `DELETE` | `/cart` | Clear entire shopping cart | Protected (`RequireAuth`) |

### 8. Orders & Checkout (`/orders`)
| Method | Endpoint | Description | Access |
| :--- | :--- | :--- | :--- |
| `POST` | `/orders/checkout` | Checkout order (atomic stock decrement & Midtrans Snap token generation) | Protected (`RequireAuth`) |
| `GET` | `/orders` | View user order history (tab filter: `all`, `unpaid`, `completed`, etc.) | Protected (`RequireAuth`) |
| `GET` | `/orders/seller` | View store orders for merchants with pagination & filters | Protected (`RequireAuth`) |
| `GET` | `/orders/:id` | Get order detail with 5-step status tracking timeline | Protected (`RequireAuth`) |
| `POST` | `/orders/:id/cancel` | Cancel an order (automatically restores product variant stock) | Protected (`RequireAuth`) |
| `PATCH` | `/orders/:id/status` | Advance order status across the fulfillment lifecycle | Protected (`RequireAuth`) |

### 9. Payments (`/payments`)
| Method | Endpoint | Description | Access |
| :--- | :--- | :--- | :--- |
| `POST` | `/payments/midtrans/webhook` | Midtrans webhook notification (SHA-512 verification, auto-restore stock on expiry) | Public |
| `GET` | `/payments/order/:order_id` | Get payment status and details for an order | Protected (`RequireAuth`) |

### 10. Health Check
| Method | Endpoint | Description | Access |
| :--- | :--- | :--- | :--- |
| `GET` | `/health` | Server health check and liveness probe | Public |

---

# Taskfile Commands

Summary of available commands via the Task runner:

| Command | Description |
| :--- | :--- |
| `task` / `task install` | Install Go dependencies and external tools (`air`, `migrate`, `swag`) |
| `task dev` | Run development server with live-reloading via Air |
| `task build` | Compile the application binary into `./tmp/main.exe` |
| `task swagger` | Generate Swagger/OpenAPI documentation from code annotations |
| `task migrate:up` | Run all pending database migrations |
| `task migrate:down` | Rollback the last applied database migration |
| `task migrate:create -- <name>` | Create a new SQL migration file with a 14-digit timestamp prefix |
| `task docker:up` | Start all services (`app`, `db`, `redis`) in the background |
| `task docker:down` | Stop and remove all Docker Compose services |
| `task docker:build` | Rebuild images and start Docker Compose services |
| `task docker:logs` | Stream live logs from all Docker Compose services |
