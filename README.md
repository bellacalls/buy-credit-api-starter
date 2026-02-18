# Sample Provider Buy Credit API

A production-ready REST API for wallet-based credit purchases with third-party partner integration.

## Architecture

Built with Clean Architecture / Domain-Driven Design (DDD) principles:

```
├── cmd/api/                    # Application entry point
├── internal/
│   ├── domain/                 # Business logic layer
│   │   ├── entity/            # Domain entities
│   │   └── repository/        # Repository interfaces
│   ├── application/           # Use cases / business logic
│   └── infrastructure/        # External concerns
│       ├── auth/              # JWT service
│       ├── http/              # HTTP handlers, middleware, response
│       └── repository/        # Repository implementations
```

## Features

- ✅ JWT-based authentication
- ✅ Wallet management
- ✅ Credit purchase transactions
- ✅ Idempotency key support
- ✅ Transaction status tracking
- ✅ Webhook registration and management
- ✅ RESTful API design
- ✅ Clean error handling
- ✅ Context-aware operations

## Quick Start

### Prerequisites
- Go 1.21+

### Installation

```bash
# Clone repository
cd /Users/dapoadeleke/GolandProjects/bella/sample-provider-buy-credit-api

# Install dependencies
go mod download

# Run the server
go run cmd/api/main.go
```

The server starts on `http://localhost:8080`

## API Documentation

### Base URL
```
http://localhost:8080/v1
```

### Authentication

**POST /auth/token**
```bash
curl -X POST http://localhost:8080/v1/auth/token \
  -H "Content-Type: application/json" \
  -d '{
    "clientId": "bella_mobile_prod",
    "clientSecret": "secret_bella_123"
  }'
```

Response:
```json
{
  "accessToken": "eyJhbGc...",
  "tokenType": "Bearer",
  "expiresIn": 3600
}
```

### Get Wallets

**GET /wallets**
```bash
curl -X GET http://localhost:8080/v1/users/{user_id}/wallets \
  -H "Authorization: Bearer {token}" \
```

Response:
```json
{
  "wallets": [
    {
      "id": "wlt_usd_abc123",
      "userId": "usr_123",
      "currency": "USD",
      "balance": "1500.50",
      "status": "ACTIVE"
    }
  ]
}
```

### Create Transaction (Buy Credit)

**POST /transactions**
```bash
curl -X POST http://localhost:8080/v1/transactions \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "idempotencyKey": "unique-key-123",
    "userId": "usr_123",
    "walletId": "wlt_usd_abc123",
    "amount": "10.00",
    "currency": "USD",
    "metadata": {
      "phoneNumber": "+1234567890",
      "provider": "bella_mobile"
    }
  }'
```

Response:
```json
{
  "transaction": {
    "id": "txn_purchase_xyz789",
    "walletId": "wlt_usd_abc123",
    "userId": "usr_123",
    "partnerWalletId": "wlt_partner_bella",
    "amount": "10.00",
    "currency": "USD",
    "status": "PENDING",
    "type": "CREDIT_PURCHASE",
    "metadata": {
      "phoneNumber": "+1234567890",
      "provider": "bella_mobile"
    },
    "createdAt": "2026-02-18T10:35:00Z",
    "updatedAt": "2026-02-18T10:35:00Z"
  }
}
```

### Get Transaction Status

**GET /transactions/{transactionId}**
```bash
curl -X GET http://localhost:8080/v1/transactions/txn_purchase_xyz789 \
  -H "Authorization: Bearer {token}"
```

Response:
```json
{
  "transaction": {
    "id": "txn_purchase_xyz789",
    "status": "SUCCESS",
    "completedAt": "2026-02-18T10:35:02Z",
    ...
  }
}
```

### Register Webhook

**POST /webhooks**
```bash
curl -X POST http://localhost:8080/v1/webhooks \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://bellamobile.co/webhooks/sample-provider",
    "events": ["transaction.completed", "transaction.failed"],
    "secret": "webhook_secret_123"
  }'
```

Response:
```json
{
  "webhook": {
    "id": "whk_abc123",
    "partnerId": "partner_bella",
    "url": "https://bellamobile.co/webhooks/sample-provider",
    "events": ["transaction.completed", "transaction.failed"],
    "status": "ACTIVE",
    "createdAt": "2026-02-18T10:00:00Z"
  }
}
```

### Get Webhooks

**GET /webhooks**
```bash
curl -X GET http://localhost:8080/v1/webhooks \
  -H "Authorization: Bearer {token}"
```

## Transaction Status Values

- `PENDING` - Transaction is being processed
- `SUCCESS` - Transaction completed successfully
- `FAILED` - Transaction failed

## Webhook Events

Webhook events are sent with HMAC signature verification:

**Headers:**
- `X-Sample Provider-Signature`: HMAC-SHA256 signature
- `X-Sample Provider-Event`: Event type

**Event Types:**
- `transaction.completed`
- `transaction.failed`

**Payload:**
```json
{
  "eventId": "evt_abc123",
  "eventType": "transaction.completed",
  "timestamp": "2026-02-18T10:35:02Z",
  "data": {
    "transaction": { ... }
  }
}
```

## Error Responses

All errors follow this format:
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message"
  }
}
```

**Common Error Codes:**
- `INVALID_CREDENTIALS` - Authentication failed
- `INVALID_TOKEN` - Token is invalid or expired
- `MISSING_AUTH_TOKEN` - No authorization header
- `INSUFFICIENT_BALANCE` - Wallet balance too low
- `WALLET_NOT_FOUND` - Wallet doesn't exist
- `TRANSACTION_NOT_FOUND` - Transaction doesn't exist

## Idempotency

Use the `Idempotency-Key` header for safe retries:
```
Idempotency-Key: unique-key-12345
```

Duplicate requests with the same key return the original transaction.

## Development

### Project Structure

- **Domain Layer**: Pure business logic, no external dependencies
- **Application Layer**: Use cases orchestrating domain logic
- **Infrastructure Layer**: HTTP handlers, repositories, external services

### Test Data

**Partner Credentials:**
- Client ID: `bella_mobile_prod`
- Client Secret: `secret_bella_123`

**Test User:**
- User ID: `usr_123`
- Wallet ID: `wlt_usd_abc123`
- Balance: `1500.50 USD`

## Production Considerations

This implementation uses in-memory storage. For production:

1. Replace repository implementations with database persistence (PostgreSQL, etc.)
2. Add proper logging and monitoring
3. Implement webhook delivery queue/retry mechanism
4. Add rate limiting
5. Implement proper secret management
6. Add comprehensive testing
7. Set up CI/CD pipeline
8. Configure HTTPS/TLS
9. Add API documentation (Swagger/OpenAPI)
10. Implement proper balance locking for transactions
