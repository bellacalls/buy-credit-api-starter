# Project Structure

## Overview

This project follows **Clean Architecture** principles with **Domain-Driven Design (DDD)** patterns, ensuring separation of concerns and maintainability.

```
sample-provider-buy-credit-api/
├── cmd/
│   └── api/
│       └── main.go                          # Application entry point
│
├── internal/
│   ├── domain/                              # Domain Layer (Pure Business Logic)
│   │   ├── entity/                          # Domain entities
│   │   │   ├── wallet.go                    # Wallet entity and value objects
│   │   │   ├── transaction.go               # Transaction entity
│   │   │   ├── webhook.go                   # Webhook entity
│   │   │   └── partner.go                   # Partner entity
│   │   │
│   │   └── repository/                      # Repository interfaces
│   │       ├── wallet_repository.go         # Wallet repository interface
│   │       ├── transaction_repository.go    # Transaction repository interface
│   │       ├── webhook_repository.go        # Webhook repository interface
│   │       └── partner_repository.go        # Partner repository interface
│   │
│   ├── application/                         # Application Layer (Use Cases)
│   │   ├── auth_usecase.go                  # Authentication business logic
│   │   ├── wallet_usecase.go                # Wallet operations logic
│   │   ├── transaction_usecase.go           # Transaction processing logic
│   │   └── webhook_usecase.go               # Webhook management logic
│   │
│   └── infrastructure/                      # Infrastructure Layer
│       ├── auth/                            # Authentication infrastructure
│       │   └── jwt_service.go               # JWT token generation/validation
│       │
│       ├── http/                            # HTTP infrastructure
│       │   ├── handler/                     # HTTP request handlers
│       │   │   ├── auth_handler.go          # Auth endpoints
│       │   │   ├── wallet_handler.go        # Wallet endpoints
│       │   │   ├── transaction_handler.go   # Transaction endpoints
│       │   │   ├── webhook_handler.go       # Webhook endpoints
│       │   │   └── router.go                # Route configuration
│       │   │
│       │   ├── middleware/                  # HTTP middleware
│       │   │   └── auth_middleware.go       # JWT authentication middleware
│       │   │
│       │   └── response/                    # HTTP response utilities
│       │       └── response.go              # Standard response formats
│       │
│       └── repository/                      # Repository implementations
│           ├── in_memory_wallet_repository.go
│           ├── in_memory_transaction_repository.go
│           ├── in_memory_webhook_repository.go
│           └── in_memory_partner_repository.go
│
├── examples/
│   └── api_test.sh                          # API testing script
│
├── docs/
│   ├── API_SPECIFICATION.md                 # Complete API documentation
│   └── PROJECT_STRUCTURE.md                 # This file
│
├── go.mod                                   # Go module definition
├── go.sum                                   # Dependency checksums
├── Makefile                                 # Build and run commands
├── README.md                                # Project documentation
└── .gitignore                               # Git ignore rules
```

---

## Layer Responsibilities

### 1. Domain Layer (`internal/domain/`)

**Purpose:** Contains pure business logic with no external dependencies.

**Components:**

- **Entities** (`entity/`):
  - Core business objects (Wallet, Transaction, Partner, Webhook)
  - Business rules and validations
  - No dependencies on frameworks or infrastructure

- **Repository Interfaces** (`repository/`):
  - Define contracts for data access
  - Keep domain layer independent of data storage details
  - Allow easy swapping of implementations

**Rules:**
- ✅ Pure Go code
- ✅ Business logic only
- ❌ No framework dependencies
- ❌ No database/HTTP/external service code

---

### 2. Application Layer (`internal/application/`)

**Purpose:** Orchestrates business logic through use cases.

**Components:**

- **Use Cases**:
  - `AuthUseCase`: Partner authentication
  - `WalletUseCase`: Wallet operations
  - `TransactionUseCase`: Transaction processing
  - `WebhookUseCase`: Webhook management

**Responsibilities:**
- Coordinate domain entities
- Implement business workflows
- Call repository interfaces
- Handle business errors

**Rules:**
- ✅ Use domain entities and repositories
- ✅ Implement business workflows
- ❌ No HTTP/database implementation details
- ❌ No framework-specific code

---

### 3. Infrastructure Layer (`internal/infrastructure/`)

**Purpose:** Handles external concerns and implementation details.

#### 3.1 Auth (`auth/`)
- JWT token generation and validation
- Cryptographic operations

#### 3.2 HTTP (`http/`)

**Handlers** (`handler/`):
- Parse HTTP requests
- Call use cases
- Format HTTP responses
- Handle HTTP-specific errors

**Middleware** (`middleware/`):
- Authentication/authorization
- Request validation
- Logging (future)
- Rate limiting (future)

**Response** (`response/`):
- Standard JSON response formatting
- Error response formatting

#### 3.3 Repository (`repository/`)
- Concrete implementations of repository interfaces
- Currently in-memory for demo
- Production: Replace with database implementations

**Rules:**
- ✅ Implement domain repository interfaces
- ✅ Framework and library usage
- ✅ External service integration
- ❌ Business logic (belongs in application layer)

---

## Dependency Flow

```
HTTP Request
    ↓
Handler (Infrastructure)
    ↓
Use Case (Application)
    ↓
Entity + Repository Interface (Domain)
    ↓
Repository Implementation (Infrastructure)
    ↓
Data Storage
```

**Key Principle:** Dependencies point inward toward the domain layer.

---

## Data Flow Examples

### Example 1: Buy Credit Transaction

```
1. Client sends POST /transactions/credit-purchase
2. Router → TransactionHandler.CreateTransaction()
3. Handler parses request, extracts auth info
4. Handler → TransactionUseCase.CreateTransaction()
5. UseCase validates via WalletRepository
6. UseCase creates Transaction entity
7. UseCase → TransactionRepository.Create()
8. UseCase processes transaction asynchronously
9. Handler formats response
10. Response sent to client
```

### Example 2: Authentication

```
1. Client sends POST /auth/token
2. Router → AuthHandler.CreateToken()
3. Handler parses credentials
4. Handler → AuthUseCase.Authenticate()
5. UseCase → PartnerRepository.FindByClientID()
6. UseCase validates credentials
7. UseCase → JWTService.GenerateToken()
8. Handler sends token response
```

---

## Design Patterns Used

### 1. **Repository Pattern**
- Abstracts data access
- Defined in domain, implemented in infrastructure
- Easy to swap implementations (in-memory → PostgreSQL)

### 2. **Use Case Pattern**
- Encapsulates business operations
- Single responsibility per use case
- Testable without infrastructure

### 3. **Dependency Injection**
- Constructor injection in `main.go`
- Enables testing with mocks
- Clear dependency graph

### 4. **Middleware Pattern**
- Cross-cutting concerns (auth, logging)
- Reusable across endpoints
- Chain of responsibility

---

## Key Design Decisions

### 1. **String-based Amounts**
```go
Balance string `json:"balance"`
```
- Avoid floating-point precision issues
- Financial calculations require exact precision
- Production: Use `decimal` library or database DECIMAL type

### 2. **Idempotency Key Support**
```go
IdempotencyKey string
```
- Prevents duplicate charges
- Stored in separate map in repository
- Production: Database with unique constraint

### 3. **Context Usage**
```go
func (r *Repository) FindByID(ctx context.Context, id string) (*Entity, error)
```
- Enables cancellation
- Timeout control
- Request tracing (future)

### 4. **Async Transaction Processing**
```go
go uc.processTransaction(context.Background(), transaction, wallet, amount)
```
- Immediate response to client
- Background processing
- Production: Use message queue (RabbitMQ, Kafka)

---

## Testing Strategy

### Unit Tests
```go
// Test domain entities
func TestTransaction_MarkSuccess(t *testing.T) { ... }

// Test use cases with mock repositories
func TestTransactionUseCase_CreateTransaction(t *testing.T) { ... }
```

### Integration Tests
```go
// Test HTTP handlers with real dependencies
func TestTransactionHandler_CreateTransaction(t *testing.T) { ... }
```

### E2E Tests
```bash
# Use examples/api_test.sh
./examples/api_test.sh
```

---

## Production Enhancements

### 1. **Database Implementation**
Replace in-memory repositories:
```go
// PostgreSQL example
type PostgresWalletRepository struct {
    db *sql.DB
}

func (r *PostgresWalletRepository) FindByID(ctx context.Context, id string) (*entity.Wallet, error) {
    var wallet entity.Wallet
    err := r.db.QueryRowContext(ctx, "SELECT * FROM wallets WHERE id = $1", id).
        Scan(&wallet.ID, &wallet.UserID, ...)
    return &wallet, err
}
```

### 2. **Message Queue**
```go
type TransactionProcessor struct {
    queue MessageQueue
}

func (uc *TransactionUseCase) CreateTransaction(...) {
    // ...
    uc.processor.Enqueue(transaction.ID)
}
```

### 3. **Observability**
```go
// Logging
log.Info("transaction created", "id", txn.ID, "amount", txn.Amount)

// Metrics
metrics.IncrementCounter("transactions_created")

// Tracing
span := trace.StartSpan(ctx, "create_transaction")
defer span.End()
```

### 4. **Configuration Management**
```go
type Config struct {
    JWTSecret     string
    DatabaseURL   string
    ServerPort    int
}

func LoadConfig() (*Config, error) {
    // Load from environment variables
}
```

---

## Adding New Features

### Example: Add Refund Endpoint

1. **Domain Layer** - Create entity/extend existing:
```go
// internal/domain/entity/transaction.go
const TransactionTypeRefund TransactionType = "REFUND"

func (t *Transaction) MarkRefunded() { ... }
```

2. **Application Layer** - Create use case:
```go
// internal/application/refund_usecase.go
func (uc *RefundUseCase) ProcessRefund(ctx context.Context, txnID string) error {
    // Business logic
}
```

3. **Infrastructure Layer** - Add handler:
```go
// internal/infrastructure/http/handler/refund_handler.go
func (h *RefundHandler) CreateRefund(w http.ResponseWriter, r *http.Request) {
    // HTTP handling
}
```

4. **Router** - Register route:
```go
// internal/infrastructure/http/handler/router.go
r.Post("/transactions/{transactionId}/refund", refundHandler.CreateRefund)
```

---

## Scaling Considerations

### Horizontal Scaling
- Stateless design enables multiple instances
- Load balancer in front
- Shared database/cache

### Database Optimization
- Indexes on frequently queried fields
- Connection pooling
- Read replicas for queries

### Caching
```go
type CachedWalletRepository struct {
    repo  repository.WalletRepository
    cache Cache
}
```

### Rate Limiting
```go
// middleware/rate_limiter.go
func RateLimitMiddleware(next http.Handler) http.Handler {
    // Implement token bucket or sliding window
}
```

---

## Security Checklist

- ✅ JWT-based authentication
- ✅ Partner credential validation
- ✅ User ID validation in transactions
- ⬜ HTTPS/TLS (production)
- ⬜ Rate limiting
- ⬜ Input validation and sanitization
- ⬜ SQL injection prevention (with prepared statements)
- ⬜ Webhook signature verification
- ⬜ Audit logging
- ⬜ Secret management (Vault, AWS Secrets Manager)

---

## References

- [Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://martinfowler.com/tags/domain%20driven%20design.html)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
