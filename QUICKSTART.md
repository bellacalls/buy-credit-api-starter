# Sample Provider Buy Credit API - Quick Start Guide

## 🚀 Get Started in 5 Minutes

### Prerequisites
- Go 1.21+ installed
- `curl` and `jq` (for testing)

### 1. Start the Server

```bash
cd /Users/dapoadeleke/GolandProjects/bella/sample-provider-buy-credit-api

# Run directly
go run cmd/api/main.go

# OR build and run
make build
./bin/api
```

Server starts on `http://localhost:8080`

### 2. Test the API

```bash
# Make test script executable
chmod +x examples/api_test.sh

# Run comprehensive test suite
./examples/api_test.sh
```

---

## 📋 Manual Testing

### Step 1: Get Access Token

```bash
curl -X POST http://localhost:8080/v1/auth/token \
  -H "Content-Type: application/json" \
  -d '{
    "clientId": "bella_mobile_prod",
    "clientSecret": "secret_bella_123"
  }'
```

**Copy the `accessToken` from the response.**

### Step 2: Get User Wallets

```bash
TOKEN="your_token_here"

curl -X GET http://localhost:8080/v1/wallets \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-User-ID: usr_123"
```

### Step 3: Buy Credit

```bash
curl -X POST http://localhost:8080/v1/transactions/credit-purchase \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-User-ID: usr_123" \
  -H "Idempotency-Key: $(uuidgen)" \
  -H "Content-Type: application/json" \
  -d '{
    "walletId": "wlt_usd_abc123",
    "amount": "10.00",
    "currency": "USD",
    "metadata": {
      "phoneNumber": "+1234567890"
    }
  }'
```

**Copy the `transaction.id` from the response.**

### Step 4: Check Transaction Status

```bash
TRANSACTION_ID="your_transaction_id_here"

curl -X GET http://localhost:8080/v1/transactions/$TRANSACTION_ID \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🏗️ Project Structure

```
sample-provider-buy-credit-api/
├── cmd/api/main.go              # Start here
├── internal/
│   ├── domain/                  # Business entities
│   ├── application/             # Use cases
│   └── infrastructure/          # HTTP, repositories
├── docs/                        # Full documentation
├── examples/api_test.sh         # Test script
└── README.md                    # Detailed guide
```

---

## 📚 Documentation

- **[README.md](README.md)** - Full project documentation
- **[API_SPECIFICATION.md](docs/API_SPECIFICATION.md)** - Complete API reference
- **[PROJECT_STRUCTURE.md](docs/PROJECT_STRUCTURE.md)** - Architecture deep dive

---

## 🧪 Test Credentials

**Partner:**
- Client ID: `bella_mobile_prod`
- Client Secret: `secret_bella_123`

**Test User:**
- User ID: `usr_123`
- Wallet ID: `wlt_usd_abc123`
- Initial Balance: `1500.50 USD`

---

## 🔑 Key Features

✅ **RESTful API** - Clean, consistent endpoints
✅ **JWT Authentication** - Secure partner access
✅ **Idempotency** - Safe transaction retries
✅ **Transaction Status** - Track purchases in real-time
✅ **Webhooks** - Event notifications
✅ **Clean Architecture** - DDD with clear separation

---

## 🛠️ Common Commands

```bash
# Run server
make run

# Build binary
make build

# Run tests
make test

# Install dependencies
make deps

# Format code
make fmt
```

---

## 🎯 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v1/auth/token` | Get access token |
| GET | `/v1/wallets` | List user wallets |
| POST | `/v1/transactions/credit-purchase` | Buy credit |
| GET | `/v1/transactions/{id}` | Get transaction status |
| POST | `/v1/webhooks` | Register webhook |
| GET | `/v1/webhooks` | List webhooks |

---

## 💡 Quick Tips

1. **Always use Idempotency-Key** for transactions to prevent duplicates
2. **Tokens expire in 1 hour** - implement refresh logic
3. **Transactions are async** - poll status or use webhooks
4. **Check balance first** to avoid insufficient balance errors

---

## 🔧 Troubleshooting

### Port Already in Use
```bash
# Find and kill process on port 8080
lsof -ti:8080 | xargs kill -9
```

### Import Errors
```bash
go mod tidy
go mod download
```

### Build Fails
```bash
# Clean and rebuild
make clean
make deps
make build
```

---

## 📞 Next Steps

1. ✅ Run the server
2. ✅ Execute `examples/api_test.sh`
3. 📖 Read [API_SPECIFICATION.md](docs/API_SPECIFICATION.md)
4. 🏗️ Explore [PROJECT_STRUCTURE.md](docs/PROJECT_STRUCTURE.md)
5. 💻 Start building your integration!

---

## 🚢 Production Checklist

Before deploying to production:

- [ ] Replace in-memory repositories with database
- [ ] Configure proper JWT secret
- [ ] Set up HTTPS/TLS
- [ ] Implement rate limiting
- [ ] Add logging and monitoring
- [ ] Set up webhook delivery queue
- [ ] Configure environment variables
- [ ] Add comprehensive tests
- [ ] Set up CI/CD pipeline
- [ ] Implement proper secret management

See [README.md](README.md) for detailed production considerations.

---

**Happy Coding! 🎉**
