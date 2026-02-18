# Architecture Overview

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Client Layer                             │
│  (Bella Mobile Widget, Third-party Apps, Mobile Apps)           │
└────────────────────────────┬────────────────────────────────────┘
                             │ HTTPS/REST
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                      API Gateway / Router                        │
│                     (chi router with middleware)                 │
│  ┌────────────────────────────────────────────────────────┐    │
│  │  Middleware: Logger | Auth | Rate Limiter | CORS       │    │
│  └────────────────────────────────────────────────────────┘    │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Infrastructure Layer                          │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │                    HTTP Handlers                          │  │
│  │  - AuthHandler                                            │  │
│  │  - WalletHandler                                          │  │
│  │  - TransactionHandler                                     │  │
│  │  - WebhookHandler                                         │  │
│  └───────────────────────┬──────────────────────────────────┘  │
└────────────────────────────┼────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Application Layer                             │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │                      Use Cases                            │  │
│  │  ┌────────────┐  ┌────────────┐  ┌──────────────┐       │  │
│  │  │   Auth     │  │  Wallet    │  │ Transaction  │       │  │
│  │  │  UseCase   │  │  UseCase   │  │   UseCase    │       │  │
│  │  └────────────┘  └────────────┘  └──────────────┘       │  │
│  │                                                            │  │
│  │  ┌────────────┐                                           │  │
│  │  │  Webhook   │                                           │  │
│  │  │  UseCase   │                                           │  │
│  │  └────────────┘                                           │  │
│  └───────────────────────┬──────────────────────────────────┘  │
└────────────────────────────┼────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Domain Layer                               │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │                  Domain Entities                          │  │
│  │  ┌──────────┐  ┌─────────────┐  ┌────────┐  ┌─────────┐ │  │
│  │  │  Wallet  │  │ Transaction │  │ Partner│  │ Webhook │ │  │
│  │  └──────────┘  └─────────────┘  └────────┘  └─────────┘ │  │
│  │                                                            │  │
│  │              Repository Interfaces (Ports)                │  │
│  │  ┌────────────────────────────────────────────────────┐  │  │
│  │  │  WalletRepo | TransactionRepo | PartnerRepo |      │  │  │
│  │  │  WebhookRepo                                        │  │  │
│  │  └────────────────────────────────────────────────────┘  │  │
│  └───────────────────────┬──────────────────────────────────┘  │
└────────────────────────────┼────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Infrastructure Layer                          │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │           Repository Implementations (Adapters)           │  │
│  │  ┌───────────────┐           ┌──────────────────┐        │  │
│  │  │   In-Memory   │    OR     │   PostgreSQL     │        │  │
│  │  │  Repositories │           │   Repositories   │        │  │
│  │  └───────────────┘           └──────────────────┘        │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
                      ┌──────────────┐
                      │   Database   │
                      │ (PostgreSQL) │
                      └──────────────┘
```

---

## Layer Communication Flow

### Request Flow Example: Buy Credit Transaction

```
1. HTTP Request
   POST /v1/transactions/credit-purchase
   Headers: Authorization, X-User-ID, Idempotency-Key
   Body: { walletId, amount, currency, metadata }
   │
   ▼
2. Router + Middleware
   - Logger: Log request
   - Auth: Validate JWT token
   - Extract user context
   │
   ▼
3. TransactionHandler (Infrastructure)
   - Parse HTTP request
   - Extract headers
   - Validate input format
   │
   ▼
4. TransactionUseCase (Application)
   - Check idempotency key
   - Validate business rules:
     * Wallet exists?
     * Wallet belongs to user?
     * Wallet active?
     * Sufficient balance?
   - Create transaction entity
   - Save via repository
   - Trigger async processing
   │
   ▼
5. Domain Entities
   - Transaction.New()
   - Wallet.ValidateBalance()
   - Business logic execution
   │
   ▼
6. Repository Interface (Domain)
   - TransactionRepository.Create()
   - WalletRepository.Update()
   │
   ▼
7. Repository Implementation (Infrastructure)
   - InMemoryTransactionRepo
   - OR PostgresTransactionRepo
   - Store in data store
   │
   ▼
8. Response
   - Success: 201 Created + Transaction JSON
   - Error: 4xx/5xx + Error JSON
```

---

## Authentication Flow

```
┌──────────┐                                    ┌──────────────┐
│  Client  │                                    │  API Server  │
└─────┬────┘                                    └──────┬───────┘
      │                                                │
      │ 1. POST /auth/token                           │
      │    { clientId, clientSecret }                 │
      ├──────────────────────────────────────────────►│
      │                                                │
      │                                      2. Validate credentials
      │                                         (PartnerRepository)
      │                                                │
      │                                      3. Generate JWT
      │                                         (JWTService)
      │                                                │
      │ 4. { accessToken, expiresIn }                 │
      │◄──────────────────────────────────────────────┤
      │                                                │
      │ 5. GET /wallets                               │
      │    Header: Authorization: Bearer {token}      │
      ├──────────────────────────────────────────────►│
      │                                                │
      │                                      6. Validate JWT
      │                                         (AuthMiddleware)
      │                                                │
      │                                      7. Extract partnerId
      │                                         from token claims
      │                                                │
      │                                      8. Execute request
      │                                                │
      │ 9. Response                                   │
      │◄──────────────────────────────────────────────┤
      │                                                │
```

---

## Transaction Processing Flow

```
┌────────┐                              ┌─────────────┐
│ Client │                              │ Transaction │
│        │                              │   UseCase   │
└───┬────┘                              └──────┬──────┘
    │                                          │
    │ 1. Create Transaction Request            │
    ├─────────────────────────────────────────►│
    │                                          │
    │                                  2. Validate Wallet
    │                                          │
    │                                  3. Check Balance
    │                                          │
    │                                  4. Create Transaction
    │                                     (Status: PENDING)
    │                                          │
    │ 5. Return Transaction (PENDING)         │
    │◄─────────────────────────────────────────┤
    │                                          │
    │                                          │
    │                              ┌───────────▼──────────┐
    │                              │  Async Processing    │
    │                              │  (Background)        │
    │                              └───────────┬──────────┘
    │                                          │
    │                              6. Deduct from customer wallet
    │                                          │
    │                              7. Credit partner wallet
    │                                          │
    │                              8. Call partner API
    │                              ┌────────────────────┐ │
    │                              │  Bella Mobile API  │ │
    │                              │  (Provision Credit)│ │
    │                              └────────────────────┘ │
    │                                          │
    │                              9. Update Transaction
    │                                 (Status: SUCCESS/FAILED)
    │                                          │
    │                              10. Send Webhook Event
    │                              ┌──────────────────┐   │
    │                              │  Partner Webhook │◄──┤
    │                              │  Endpoint        │
    │                              └──────────────────┘
    │                                          │
    │ 11. Poll Transaction Status              │
    ├─────────────────────────────────────────►│
    │                                          │
    │ 12. Return Transaction (SUCCESS)         │
    │◄─────────────────────────────────────────┤
    │                                          │
```

---

## Webhook Event Delivery

```
┌─────────────────┐                         ┌──────────────────┐
│ Transaction     │                         │ Partner Server   │
│ Processing      │                         │ (Bella Mobile)   │
└────────┬────────┘                         └────────┬─────────┘
         │                                           │
         │ 1. Transaction completes                 │
         │    (SUCCESS or FAILED)                   │
         │                                           │
         │ 2. Create webhook event                  │
         ├──────────────┐                           │
         │              │                           │
         │ 3. Find registered webhooks              │
         │    for partner                           │
         │              │                           │
         │ 4. Generate HMAC signature               │
         │    (using webhook secret)                │
         │              │                           │
         │ 5. POST to webhook URL                   │
         │    Headers:                              │
         │      X-Sample Provider-Signature: sha256=...        │
         │      X-Sample Provider-Event: transaction.completed │
         │    Body:                                 │
         │      { eventId, eventType, data }        │
         ├─────────────────────────────────────────►│
         │                                           │
         │                              6. Verify signature
         │                                 (HMAC validation)
         │                                           │
         │                              7. Process event
         │                                           │
         │ 8. 200 OK                                │
         │◄─────────────────────────────────────────┤
         │                                           │
         │ [If failed: retry with backoff]          │
         │ Retry schedule: 1m, 5m, 15m, 1h, 6h      │
         │                                           │
```

---

## Data Model Relationships

```
┌──────────────┐
│   Partner    │
│──────────────│
│ id           │───┐
│ clientId     │   │
│ clientSecret │   │
│ walletId     │───┼───────┐
│ status       │   │       │
└──────────────┘   │       │
                   │       │
                   │       │ Has Wallet
                   │       │
┌──────────────┐   │       │        ┌──────────────┐
│    Wallet    │◄──┼───────┘   ┌────│ Transaction  │
│──────────────│   │            │    │──────────────│
│ id           │   │            │    │ id           │
│ userId       │   │            │    │ walletId     │───┐
│ currency     │   │            │    │ userId       │   │
│ balance      │   │            │    │ partnerWalletId│─┤
│ status       │   │            │    │ amount       │   │
└──────┬───────┘   │            │    │ currency     │   │
       │           │            │    │ status       │   │
       │           │            │    │ type         │   │
       │ Has Many  │            │    │ metadata     │   │
       │           │            │    └──────────────┘   │
       └───────────┼────────────┘                       │
                   │                                    │
                   │                                    │
                   │ Has Many                           │
                   │                                    │
┌──────────────┐   │                                    │
│   Webhook    │◄──┘                                    │
│──────────────│                         References     │
│ id           │                         (both ways)    │
│ partnerId    │◄────────────────────────────────────────┘
│ url          │
│ events[]     │
│ secret       │
│ status       │
└──────────────┘

Events Flow:
- Partner registers webhook
- Transaction created → triggers event
- Event sent to webhook URL
- Partner processes notification
```

---

## Clean Architecture Boundaries

```
┌───────────────────────────────────────────────────────────────┐
│                         External World                         │
│  (HTTP, Database, Message Queues, External APIs)              │
└───────────────────────────┬───────────────────────────────────┘
                            │
                ┌───────────▼──────────┐
                │  Infrastructure      │
                │  (Adapters)          │
                │  - Handlers          │
                │  - Repositories      │
                │  - Auth Service      │
                └───────────┬──────────┘
                            │
                    Dependencies point
                        inward →
                            │
                ┌───────────▼──────────┐
                │  Application         │
                │  (Use Cases)         │
                │  - Business Logic    │
                │  - Orchestration     │
                └───────────┬──────────┘
                            │
                            │
                ┌───────────▼──────────┐
                │  Domain              │
                │  (Entities + Ports)  │
                │  - Business Rules    │
                │  - Interfaces        │
                │  - No Dependencies   │
                └──────────────────────┘

Rules:
1. Domain doesn't know about Application or Infrastructure
2. Application depends on Domain, not Infrastructure
3. Infrastructure depends on both Domain and Application
4. Dependencies only point inward (towards Domain)
```

---

## Dependency Injection

```go
func main() {
    // Infrastructure Layer
    walletRepo := repository.NewInMemoryWalletRepository()
    transactionRepo := repository.NewInMemoryTransactionRepository()
    webhookRepo := repository.NewInMemoryWebhookRepository()
    partnerRepo := repository.NewInMemoryPartnerRepository()

    jwtService := auth.NewJWTService("secret-key")

    // Application Layer (inject repositories)
    authUseCase := application.NewAuthUseCase(partnerRepo, jwtService)
    walletUseCase := application.NewWalletUseCase(walletRepo)
    transactionUseCase := application.NewTransactionUseCase(
        transactionRepo,
        walletRepo,
        webhookRepo,
    )
    webhookUseCase := application.NewWebhookUseCase(webhookRepo)

    // Infrastructure Layer (inject use cases)
    authHandler := handler.NewAuthHandler(authUseCase)
    walletHandler := handler.NewWalletHandler(walletUseCase)
    transactionHandler := handler.NewTransactionHandler(transactionUseCase)
    webhookHandler := handler.NewWebhookHandler(webhookUseCase)

    // Wire everything together
    router := handler.SetupRouter(
        authHandler,
        walletHandler,
        transactionHandler,
        webhookHandler,
        authMiddleware,
    )
}
```

---

## Scaling Architecture

### Horizontal Scaling

```
                     ┌────────────────┐
                     │  Load Balancer │
                     └────────┬───────┘
                              │
         ┌────────────────────┼────────────────────┐
         │                    │                    │
         ▼                    ▼                    ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   API Server    │  │   API Server    │  │   API Server    │
│   Instance 1    │  │   Instance 2    │  │   Instance 3    │
└────────┬────────┘  └────────┬────────┘  └────────┬────────┘
         │                    │                    │
         └────────────────────┼────────────────────┘
                              │
                   ┌──────────┴──────────┐
                   │                     │
                   ▼                     ▼
         ┌──────────────────┐   ┌──────────────────┐
         │    PostgreSQL    │   │   Redis Cache    │
         │   (Primary)      │   │                  │
         └─────────┬────────┘   └──────────────────┘
                   │
                   ▼
         ┌──────────────────┐
         │    PostgreSQL    │
         │   (Read Replica) │
         └──────────────────┘
```

### Message Queue Integration

```
┌─────────────────┐        ┌──────────────────┐
│  API Server     │        │ Message Queue    │
│                 │        │ (RabbitMQ/Kafka) │
└────────┬────────┘        └─────────┬────────┘
         │                           │
         │ 1. Create Transaction     │
         │    (Status: PENDING)      │
         │                           │
         │ 2. Enqueue Job            │
         ├──────────────────────────►│
         │                           │
         │ 3. Return to Client       │
         │    (202 Accepted)         │
         │                           │
         │                           │
         │                  ┌────────▼─────────┐
         │                  │  Worker Pool     │
         │                  │  (Processors)    │
         │                  └────────┬─────────┘
         │                           │
         │                  4. Dequeue & Process
         │                           │
         │                  5. Update Transaction
         │                           │
         │                  6. Send Webhook
         │                           │
```

---

## Security Architecture

```
┌───────────────────────────────────────────────────────┐
│                   Security Layers                      │
├───────────────────────────────────────────────────────┤
│                                                        │
│  1. Transport Layer Security (TLS/HTTPS)              │
│     └─► Encrypt data in transit                       │
│                                                        │
│  2. API Gateway / WAF                                 │
│     ├─► DDoS protection                               │
│     ├─► IP whitelisting                               │
│     └─► Rate limiting                                 │
│                                                        │
│  3. Authentication (JWT)                              │
│     ├─► Partner verification                          │
│     ├─► Token validation                              │
│     └─► Expiration checks                             │
│                                                        │
│  4. Authorization                                     │
│     ├─► User ownership validation                     │
│     ├─► Resource access control                       │
│     └─► Role-based permissions                        │
│                                                        │
│  5. Input Validation                                  │
│     ├─► Schema validation                             │
│     ├─► SQL injection prevention                      │
│     └─► XSS protection                                │
│                                                        │
│  6. Audit Logging                                     │
│     ├─► Request/response logging                      │
│     ├─► Transaction audit trail                       │
│     └─► Security event monitoring                     │
│                                                        │
│  7. Data Encryption                                   │
│     ├─► Encrypt sensitive data at rest               │
│     ├─► Secret management (Vault)                     │
│     └─► Key rotation                                  │
│                                                        │
└───────────────────────────────────────────────────────┘
```

---

## Observability Stack

```
┌────────────────────────────────────────────────────────┐
│                   API Servers                           │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Instrumentation:                                 │  │
│  │  - Structured logging                             │  │
│  │  - Metrics collection                             │  │
│  │  - Distributed tracing                            │  │
│  └──────────────────────────────────────────────────┘  │
└───────┬────────────────┬────────────────┬──────────────┘
        │                │                │
        │                │                │
        ▼                ▼                ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│   Logging    │ │   Metrics    │ │   Tracing    │
│   (ELK/Loki) │ │ (Prometheus) │ │   (Jaeger)   │
└──────┬───────┘ └──────┬───────┘ └──────┬───────┘
       │                │                │
       └────────────────┼────────────────┘
                        │
                        ▼
               ┌──────────────────┐
               │    Dashboards    │
               │   (Grafana)      │
               └──────────────────┘
                        │
                        ▼
               ┌──────────────────┐
               │    Alerting      │
               │  (PagerDuty)     │
               └──────────────────┘

Metrics to Track:
- Request rate (requests/sec)
- Error rate (errors/sec)
- Response time (p50, p95, p99)
- Transaction success rate
- Wallet balance changes
- Webhook delivery success
```

---

This architecture ensures:

✅ **Separation of Concerns** - Clear layer boundaries
✅ **Testability** - Easy to mock dependencies
✅ **Maintainability** - Easy to understand and modify
✅ **Scalability** - Stateless design enables horizontal scaling
✅ **Security** - Multiple security layers
✅ **Observability** - Comprehensive monitoring and logging
