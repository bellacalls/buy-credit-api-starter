# Buy Credit API Quickstart

A production-ready REST API for credit purchases.

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
- ✅ Credit purchase transactions
- ✅ Transaction status tracking
- ✅ RESTful API design
- ✅ Clean error handling
- ✅ Context-aware operations

## Quick Start

### Prerequisites
- Go 1.21+

### Installation

```bash
# Clone repository
cd buy-credit-api-starter

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

### 1. Get Access Token

**POST /auth/token**

Request:
```bash
curl -X POST http://localhost:8080/v1/auth/token \
  -H "Content-Type: application/json" \
  -d '{
    "apiKey": "API_KEY",
    "apiSecret": "API_SECRET"
  }'
```

Response:
```json
{
  "accessToken": "ACCESS_TOKEN"
}
```

### 2. Transaction (Buy Credit)

**POST /transactions**

Request:
```bash
curl -X POST http://localhost:8080/v1/transactions \
  -H "Authorization: Bearer ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "user_123",
    "amount": 10,
    "currency": "USD"
  }'
```

Response:
```json
{
  "transactionId": "txn_123",
  "userId": "user_123",
  "currency": "USD",
  "amount": 10,
  "status": "SUCCESSFUL",
  "timestamp": "2026-02-05T10:35:00Z"
}
```

### 3. Get Transaction (Buy Credit Status)

**GET /transactions/{transactionId}**

Request:
```bash
curl -X GET http://localhost:8080/v1/transactions/txn_123 \
  -H "Authorization: Bearer ACCESS_TOKEN"
```

Response:
```json
{
  "transactionId": "txn_123",
  "userId": "user_123",
  "currency": "USD",
  "amount": 10,
  "status": "SUCCESSFUL",
  "timestamp": "2026-02-05T10:35:00Z"
}
```

## Transaction Status Values

- `PENDING` - Transaction is being processed
- `SUCCESSFUL` - Transaction completed successfully
- `FAILED` - Transaction failed

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
- `INVALID_AMOUNT` - Amount is invalid or negative
- `TRANSACTION_NOT_FOUND` - Transaction doesn't exist

## Development

### Project Structure

- **Domain Layer**: Pure business logic, no external dependencies
- **Application Layer**: Use cases orchestrating domain logic
- **Infrastructure Layer**: HTTP handlers, repositories, external services

### Test Data

**API Credentials:**
- API Key: `bella_mobile_prod`
- API Secret: `secret_bella_123`

**Test User:**
- User ID: `user_123`

## Production Considerations

This implementation uses in-memory storage. For production:

1. **Database Persistence**: Replace repository implementations with database (PostgreSQL, etc.)
2. **Logging & Monitoring**: Add structured logging and APM
3. **Rate Limiting**: Implement rate limiting on all endpoints
4. **Secret Management**: Use secure secret storage (AWS Secrets Manager, Vault, etc.)
5. **Testing**: Add comprehensive unit, integration, and security tests
6. **CI/CD**: Set up automated testing and deployment pipeline
7. **HTTPS/TLS**: Configure TLS certificates and enforce HTTPS
8. **API Documentation**: Generate Swagger/OpenAPI documentation
9. **Transaction Processing**: Add proper transaction state machine and timeout handling
10. **Webhook Notifications**: Implement webhook delivery for transaction status updates
