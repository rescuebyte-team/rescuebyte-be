# RescueByte Backend - Architectural Guide & Standards

This documentation serves as the architectural reference and engineering standards for the **RescueByte (`rescuebyte-be`)** backend project.

## 📌 Core Technologies
- **Language:** Go (v1.25+)
- **HTTP Web Framework:** GoFiber v3 (`github.com/gofiber/fiber/v3`)
- **Database Driver & Query:** PostgreSQL (`pgx/v5`, `sqlx`, `Masterminds/squirrel`)
- **Cache & Key-Value:** Redis (`go-redis/v9`)
- **Authentication & Authorization:** JWT (`golang-jwt/jwt/v5`), Google OAuth2 (`golang.org/x/oauth2`), Bcrypt (`golang.org/x/crypto`), Role-Based Access Control (RBAC)
- **Payment Gateway:** Midtrans Snap API Sandbox (`pkg/midtrans`) with SHA-512 signature verification
- **Cloud Storage:** Supabase Storage (`pkg/supabase`)
- **Task Runner:** Task (`Taskfile.yml`)
- **API Documentation:** Swagger / OpenAPI 2.0 (`swaggo/swag`, `gofiber/contrib/v3/swaggo`)

---

## 🗂️ Directory Structure

```text
rescuebyte-be/
├── .gemini/skills/rescuebyte-architecture/   # Custom skill for AI pair-programming
├── cmd/
│   └── app/
│       └── main.go                           # HTTP server application entrypoint
├── config/
│   ├── .env.example                          # Environment configuration template
│   └── .env                                  # Local environment configuration file
├── database/
│   └── migrations/                           # PostgreSQL database migrations (golang-migrate format)
├── domain/                                   # Domain Core (framework-agnostic business layer)
│   ├── contracts/                            # Repository & Service interfaces (Dependency Inversion)
│   ├── dto/                                  # ALL Request, Query, Response, & Integration DTOs
│   ├── entity/                               # Database table entities (DB column mappings)
│   ├── enums/                                # Rich Enums (UserRole, OrderStatus, PaymentMethod, etc.)
│   ├── errx/                                 # Centralized structured error types (RequestError)
│   └── helpers/                              # Domain helpers (pagination, etc.)
├── docs/                                     # Swagger documentation & architecture guide
│   ├── ARCHITECTURE_GUIDE.md                 # System architectural guide & conventions
│   ├── docs.go                               # Auto-generated Swagger spec Go code
│   ├── swagger.json                          # Auto-generated Swagger JSON
│   └── swagger.yaml                          # Auto-generated Swagger YAML
├── internal/
│   ├── app/                                  # Feature-Sliced Application Modules
│   │   ├── auth/                             # Authentication & Registration Module
│   │   ├── base/                             # Base controller / generic handler
│   │   ├── cart/                             # Marketplace Shopping Cart Module
│   │   ├── enum/                             # System Enums Lookup Module
│   │   ├── order/                            # Checkout, Orders, & Status Tracking Module
│   │   ├── payment/                          # Payments & Midtrans Webhook Module
│   │   ├── product/                          # Mystery Box Products & Marketplace Cards Module
│   │   ├── storage/                          # File Upload to Supabase Storage Module
│   │   ├── store/                            # Store Profiles & Partner Operational Hours Module
│   │   └── user/                             # User Profile & RBAC Module
│   ├── infra/
│   │   ├── database/                         # Connection pool, Transactor, DBGetter
│   │   ├── env/                              # Viper environment configuration loader
│   │   ├── redis/                            # Redis client connection
│   │   └── server/                           # Fiber v3 HTTP server & route mounting
│   └── middlewares/                          # JWT Auth, RBAC, RequestID, Logger, Recover, CORS
├── pkg/                                      # Reusable Packages & Clients (Must be Domain-Agnostic)
│   ├── helpers/http/error_handler/           # Fiber centralized error handler
│   ├── helpers/http/response/                # Standardized JSON response wrapper
│   ├── jwt/                                  # JWT token generator & claims verifier
│   ├── log/                                  # Zerolog structured logging wrapper
│   ├── midtrans/                             # Midtrans Snap Client & SHA-512 Verifier (Uses DTOs from domain/dto)
│   ├── oauth/                                # Google OAuth2 integration client
│   ├── supabase/                             # Supabase Storage client
│   ├── uuid/                                 # UUID v4 and v7 helper
│   └── validator/                            # Go-playground request validator
├── Taskfile.yml                              # Task automation commands (dev, build, migrate, swagger)
└── go.mod
```

---

## 🏛️ Mandatory Architectural Principles & Conventions

### 1. Placement of Data Transfer Objects (DTO)
> [!IMPORTANT]
> **ALL DTO STRUCTS MUST RESIDE IN `domain/dto/`!**
> The `pkg/` directory (e.g., `pkg/midtrans`) contains domain-agnostic, reusable utility and client packages. Under no circumstances should business DTOs or webhook payloads be declared inside `pkg/`. All request payloads, responses, and external integration contracts (such as `SnapTransactionRequest`, `SnapTransactionResponse`, and `MidtransNotificationPayload`) must be defined in `domain/dto/payment_dto.go`.

### 2. File Separation Convention per Layer (`internal/app/<domain>/`)
Each feature module cleanly decouples struct definitions/constructors from method implementations:
- **`repository/`**: `repository.go` (struct definition & `New...` constructor), `<domain>_repository.go` (Squirrel database queries).
- **`service/`**: `service.go` (struct definition & `New...` constructor), `<domain>_service.go` (core business logic).
- **`controller/`**: `controller.go` (struct definition & route mounting via `Init<Domain>Controller`), `<domain>_controller.go` (HTTP handlers).

### 3. Atomic Database Transactions (`Transactor` & `DBGetter`)
Repositories do not initiate their own database transactions; instead, they receive a `database.DBGetter`. Any multi-table mutation (for example, checkout operations that deduct variant inventory and clear cart items) must be wrapped inside `s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error { ... })`.

### 4. Atomic Inventory / Stock Management
Inventory deductions for product variants must execute via conditional SQL updates to prevent race conditions:
```sql
UPDATE product_variants
SET current_stock = current_stock - $1, updated_at = NOW()
WHERE id = $2 AND current_stock >= $1;
```
If `rowsAffected == 0`, the system returns `errx.ErrInsufficientStock` and aborts the entire transaction.

---

## 🚀 RescueByte API Endpoint List (v1)

Base URL: `http://localhost:8080/api/v1`

### 1. Authentication (`/auth`)
| Method | Endpoint | Description | Protection |
| :--- | :--- | :--- | :--- |
| `POST` | `/auth/register` | Register a new user account (email, password, name) | Public |
| `POST` | `/auth/login` | Login using email and password | Public |
| `POST` | `/auth/refresh` | Renew Access Token using Refresh Token | Public |
| `POST` | `/auth/logout` | Revoke Refresh Token and sign out | Public |
| `GET` | `/auth/google` | Get Google OAuth2 consent redirect URL | Public |
| `GET` | `/auth/google/callback` | Google OAuth2 callback (automatic sign-in / registration) | Public |
| `GET` | `/auth/me` | Fetch authenticated user's active session profile | `RequireAuth()` |

### 2. Users (`/users`)
| Method | Endpoint | Description | Protection |
| :--- | :--- | :--- | :--- |
| `GET` | `/users/profile` | Retrieve profile data for the authenticated user | `RequireAuth()` |
| `PUT` | `/users/profile` | Update user profile (name, avatar, phone, gender, date of birth) | `RequireAuth()` |
| `PATCH` | `/users/role` | Update user role | `RequireAuth()`, `RequireRoles(RoleAdminIdx)` |

### 3. System Enums Lookup (`/enums`)
| Method | Endpoint | Description | Protection |
| :--- | :--- | :--- | :--- |
| `GET` | `/enums` | Retrieve system enum dictionary (supports optional query filter `?types=user_role,order_status`) | Public |

### 4. Cloud Storage (`/storage`)
| Method | Endpoint | Description | Protection |
| :--- | :--- | :--- | :--- |
| `POST` | `/storage/upload` | Upload a single file (image/document) to Supabase Storage | `RequireAuth()` |
| `POST` | `/storage/upload-multiple` | Upload multiple files simultaneously to Supabase Storage | `RequireAuth()` |

### 5. Partner Stores (`/stores`)
| Method | Endpoint | Description | Protection |
| :--- | :--- | :--- | :--- |
| `GET` | `/stores` | List registered stores (filter by `name`, `city`, geolocation radius `latitude` & `longitude`, pagination) | Public |
| `GET` | `/stores/:id` | Get store details by UUID | Public |
| `GET` | `/stores/:id/summary` | Get average rating and total reviews for a store | Public |
| `POST` | `/stores` | Register a new store for a merchant | `RequireAuth()` |
| `GET` | `/stores/me/profile` | Retrieve store profile of the authenticated user | `RequireAuth()` |
| `PUT` | `/stores/me/profile` | Update store information (name, description, address, profile picture) | `RequireAuth()` |
| `PUT` | `/stores/me/operational` | Update store operational hours and working days | `RequireAuth()` |
| `PATCH` | `/stores/me/status` | Toggle store operational status (open / closed) | `RequireAuth()` |

### 6. Products & Marketplace Cards (`/products`)
| Method | Endpoint | Description | Protection |
| :--- | :--- | :--- | :--- |
| `GET` | `/products` | Retrieve lightweight mystery box product cards for Marketplace UI (joins store & variant data; filters: `search`, `city`, `latitude`, `longitude`, `page`, `limit`) | Public |
| `GET` | `/products/:id` | Full product details including price variants and image gallery | Public |
| `POST` | `/products` | Create a new mystery box product with variants and images | `RequireAuth()` |
| `GET` | `/products/me/list` | List products owned by the authenticated merchant's store | `RequireAuth()` |
| `PUT` | `/products/:id` | Update product details | `RequireAuth()` |
| `PATCH` | `/products/:id/status` | Activate / deactivate a product | `RequireAuth()` |
| `DELETE` | `/products/:id` | Soft delete a product | `RequireAuth()` |

### 7. Shopping Cart (`/cart`)
| Method | Endpoint | Description | Protection |
| :--- | :--- | :--- | :--- |
| `GET` | `/cart` | Retrieve all cart items grouped per store, complete with discount calculations, subtotal, and distance (query: `latitude`, `longitude`) | `RequireAuth()` |
| `GET` | `/cart/count` | Retrieve total quantity of items in the user's cart | `RequireAuth()` |
| `POST` | `/cart/items` | Add a product variant to the cart | `RequireAuth()` |
| `PATCH` | `/cart/items/:id` | Update quantity of a specific cart item | `RequireAuth()` |
| `DELETE` | `/cart/items/:id` | Remove a single item from the cart | `RequireAuth()` |
| `DELETE` | `/cart/stores/:store_id` | Remove all items belonging to a specific store from the cart | `RequireAuth()` |
| `DELETE` | `/cart` | Empty the entire shopping cart | `RequireAuth()` |

### 8. Orders (`/orders`)
| Method | Endpoint | Description | Protection |
| :--- | :--- | :--- | :--- |
| `POST` | `/orders/checkout` | Checkout order (atomic stock deduction, creates Midtrans Snap transaction token for online payment) | `RequireAuth()` |
| `GET` | `/orders` | Retrieve user order history (tab filter: `all`, `unpaid`, `waiting_confirmation`, `ready_pickup`, `completed`, `cancelled`, pagination) | `RequireAuth()` |
| `GET` | `/orders/seller` | Retrieve incoming orders for the authenticated merchant's store, with pagination (`page`, `limit`), status filter (`tab`), and search (`search`) | `RequireAuth()` |
| `GET` | `/orders/:id` | Retrieve full order details including a 5-step status tracking timeline, cost breakdown, and store info | `RequireAuth()` |
| `POST` | `/orders/:id/cancel` | Cancel an order (automatically restores product variant stock) | `RequireAuth()` |
| `PATCH` | `/orders/:id/status` | Update order status throughout the order lifecycle | `RequireAuth()` |

*Example seller orders query:* `GET /api/v1/orders/seller?page=1&limit=10&tab=WAITING_PAYMENT&search=RB`. The endpoint determines the store from the authenticated account, so clients do not pass `store_id`.

### 9. Payments (`/payments`)
| Method | Endpoint | Description | Protection |
| :--- | :--- | :--- | :--- |
| `POST` | `/payments/midtrans/webhook` | Webhook notification from Midtrans (verifies SHA-512 signature, transitions payment and order status, auto-restores inventory upon expiration/cancellation) | Public |
| `GET` | `/payments/order/:order_id` | Check payment status and transaction details for an order | `RequireAuth()` |

### 10. Health Check
| Method | Endpoint | Description | Protection |
| :--- | :--- | :--- | :--- |
| `GET` | `/health` | Check API server availability and liveness | Public |

---

## 📖 Swagger API Documentation

RescueByte Backend comes equipped with interactive **Swagger UI / OpenAPI 2.0** documentation:
- **Swagger UI Access:** `http://localhost:8080/swagger/index.html` (or `http://localhost:8080/swagger`)
- **JSON Specification:** `http://localhost:8080/docs/swagger.json`
- **YAML Specification:** `http://localhost:8080/docs/swagger.yaml`
- **Authorizing in Swagger:** Click the **Authorize** button in the top-right corner of Swagger UI and enter your JWT token formatted as:
  ```text
  Bearer <access_token>
  ```
- **Regenerating Documentation:** Whenever adding or updating endpoints/DTOs, run:
  ```bash
  task swagger
  # or
  swag init -g ./cmd/app/main.go -o ./docs --parseDependency --parseInternal
  ```

---

## 🛠️ Quick Commands

```bash
# Run application with live-reloading (Air)
task dev

# Generate Swagger API documentation
task swagger

# Build application binary
task build

# Run database migrations (PostgreSQL)
task migrate:up

# Rollback 1 migration step
task migrate:down

# Start the full stack via Docker Compose
task docker:up
# or: docker compose up -d

# Stop Docker Compose services
task docker:down
# or: docker compose down
```
