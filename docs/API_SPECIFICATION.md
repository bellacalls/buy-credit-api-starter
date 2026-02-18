# Buy Credit API Specification Quickstart

## Overview

This API enables third-party partners (like Bella Mobile) to sell airtime and prepaid credits to sample-provider customers using their sample-provider wallet balance.

**Base URL:** `https://api.sample-provider.co/v1`

---

## Authentication

### Get Access Token

Obtain a JWT token for API access.

**Endpoint:** `POST /auth/token`

**Request:**
```json
{
  "clientId": "bella_mobile_prod",
  "clientSecret": "secret_xyz"
}
```

**Success Response (200 OK):**
```json
{
  "accessToken": "eyJhbGc...",
  "tokenType": "Bearer",
  "expiresIn": 3600
}
```

**Error Response (401 Unauthorized):**
```json
{
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Invalid client credentials"
  }
}
```

**Error Codes:**
- `INVALID_CREDENTIALS` - Wrong clientId or clientSecret
- `MISSING_FIELDS` - Required fields are missing

---

## Wallet Operations

### Get User Wallets

Retrieve all wallets for a specific user.

**Endpoint:** `GET /user/{user_id}/wallets`

**Headers:**
```
Authorization: Bearer {accessToken}
X-User-ID: {userId}
```

**Success Response (200 OK):**
```json
{
  "wallets": [
    {
      "id": "wlt_usd_abc123",
      "userId": "usr_123",
      "currency": "USD",
      "balance": "1500.50",
      "status": "ACTIVE",
      "createdAt": "2026-02-18T10:00:00Z",
      "updatedAt": "2026-02-18T10:00:00Z"
    }
  ]
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": {
    "code": "MISSING_USER_ID",
    "message": "X-User-ID header is required"
  }
}
```

**Wallet Status Values:**
- `ACTIVE` - Wallet is operational
- `INACTIVE` - Wallet is disabled
- `FROZEN` - Wallet is temporarily suspended

---

## Transaction Operations

### Create Transaction (Buy Credit)

Purchase credit/airtime by transferring funds from customer wallet to partner wallet.

**Endpoint:** `POST /transactions`

**Headers:**
```
Authorization: Bearer {accessToken}
X-User-ID: {userId}
Idempotency-Key: {unique_key} (optional but recommended)
Content-Type: application/json
```

**Request:**
```json
{
  "walletId": "wlt_usd_abc123",
  "amount": "10.00",
  "currency": "USD",
  "metadata": {
    "phoneNumber": "+1234567890",
    "provider": "bella_mobile",
    "productId": "airtime_10usd"
  }
}
```

**Success Response (201 Created):**
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
      "provider": "bella_mobile",
      "productId": "airtime_10usd"
    },
    "createdAt": "2026-02-18T10:35:00Z",
    "updatedAt": "2026-02-18T10:35:00Z"
  }
}
```

**Error Responses:**

*400 Bad Request - Insufficient Balance:*
```json
{
  "error": {
    "code": "INSUFFICIENT_BALANCE",
    "message": "Wallet balance is insufficient"
  }
}
```

*404 Not Found - Wallet Not Found:*
```json
{
  "error": {
    "code": "WALLET_NOT_FOUND",
    "message": "wallet not found"
  }
}
```

*403 Forbidden - Wrong Wallet Owner:*
```json
{
  "error": {
    "code": "FORBIDDEN",
    "message": "wallet does not belong to user"
  }
}
```

**Transaction Status Values:**
- `PENDING` - Transaction is being processed
- `SUCCESS` - Transaction completed successfully
- `FAILED` - Transaction failed

**Transaction Types:**
- `CREDIT_PURCHASE` - Purchase of airtime/prepaid credit

**Idempotency:**
- Include `Idempotency-Key` header to prevent duplicate charges
- Repeated requests with the same key return the original transaction
- Keys should be unique per transaction attempt
- Example: UUID, timestamp-based string, or request hash

---

### Get Transaction Status

Query the status of a specific transaction.

**Endpoint:** `GET /transactions/{transactionId}`

**Headers:**
```
Authorization: Bearer {accessToken}
```

**Success Response (200 OK):**
```json
{
  "transaction": {
    "id": "txn_purchase_xyz789",
    "walletId": "wlt_usd_abc123",
    "userId": "usr_123",
    "partnerWalletId": "wlt_partner_bella",
    "amount": "10.00",
    "currency": "USD",
    "status": "SUCCESS",
    "type": "CREDIT_PURCHASE",
    "metadata": {...},
    "createdAt": "2026-02-18T10:35:00Z",
    "updatedAt": "2026-02-18T10:35:02Z",
    "completedAt": "2026-02-18T10:35:02Z"
  }
}
```

**Error Response (404 Not Found):**
```json
{
  "error": {
    "code": "TRANSACTION_NOT_FOUND",
    "message": "Transaction not found"
  }
}
```

---

## Webhook Management

### Register Webhook

Register a webhook endpoint to receive transaction status notifications.

**Endpoint:** `POST /webhooks`

**Headers:**
```
Authorization: Bearer {accessToken}
Content-Type: application/json
```

**Request:**
```json
{
  "url": "https://bellamobile.co/webhooks/sample-provider",
  "events": ["transaction.completed", "transaction.failed"],
  "secret": "webhook_secret_123"
}
```

**Success Response (201 Created):**
```json
{
  "webhook": {
    "id": "whk_abc123",
    "partnerId": "partner_bella",
    "url": "https://bellamobile.co/webhooks/sample-provider",
    "events": ["transaction.completed", "transaction.failed"],
    "status": "ACTIVE",
    "createdAt": "2026-02-18T10:00:00Z",
    "updatedAt": "2026-02-18T10:00:00Z"
  }
}
```

**Supported Events:**
- `transaction.completed` - Transaction finished successfully
- `transaction.failed` - Transaction failed

---

### Get Registered Webhooks

List all webhooks registered by the authenticated partner.

**Endpoint:** `GET /webhooks`

**Headers:**
```
Authorization: Bearer {accessToken}
```

**Success Response (200 OK):**
```json
{
  "webhooks": [
    {
      "id": "whk_abc123",
      "partnerId": "partner_bella",
      "url": "https://bellamobile.co/webhooks/sample-provider",
      "events": ["transaction.completed", "transaction.failed"],
      "status": "ACTIVE",
      "createdAt": "2026-02-18T10:00:00Z",
      "updatedAt": "2026-02-18T10:00:00Z"
    }
  ]
}
```

---

## Webhook Event Format

When a registered event occurs, sample-provider sends an HTTP POST request to your webhook URL.

**Headers:**
```
X-sample-provider-Signature: sha256=<hmac_signature>
X-sample-provider-Event: transaction.completed
Content-Type: application/json
```

**Payload:**
```json
{
  "eventId": "evt_abc123",
  "eventType": "transaction.completed",
  "timestamp": "2026-02-18T10:35:02Z",
  "data": {
    "transaction": {
      "id": "txn_purchase_xyz789",
      "walletId": "wlt_usd_abc123",
      "userId": "usr_123",
      "partnerWalletId": "wlt_partner_bella",
      "amount": "10.00",
      "currency": "USD",
      "status": "SUCCESS",
      "type": "CREDIT_PURCHASE",
      "metadata": {...},
      "createdAt": "2026-02-18T10:35:00Z",
      "completedAt": "2026-02-18T10:35:02Z"
    }
  }
}
```

**Signature Verification:**
```python
import hmac
import hashlib

def verify_webhook(payload, signature, secret):
    expected_signature = hmac.new(
        secret.encode('utf-8'),
        payload.encode('utf-8'),
        hashlib.sha256
    ).hexdigest()

    return hmac.compare_digest(
        f"sha256={expected_signature}",
        signature
    )
```

**Expected Response:**
- Return `200 OK` to acknowledge receipt
- sample-provider will retry failed deliveries (non-200 responses)
- Retry schedule: 1min, 5min, 15min, 1hr, 6hr

---

## Error Response Format

All errors follow this consistent structure:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message"
  }
}
```

### Common Error Codes

| Code | Status | Description |
|------|--------|-------------|
| `INVALID_CREDENTIALS` | 401 | Authentication failed |
| `INVALID_TOKEN` | 401 | Token is invalid or expired |
| `MISSING_AUTH_TOKEN` | 401 | No authorization header |
| `UNAUTHORIZED` | 401 | Invalid authentication |
| `FORBIDDEN` | 403 | Access denied |
| `WALLET_NOT_FOUND` | 404 | Wallet doesn't exist |
| `TRANSACTION_NOT_FOUND` | 404 | Transaction doesn't exist |
| `INVALID_REQUEST` | 400 | Malformed request body |
| `MISSING_FIELDS` | 400 | Required fields missing |
| `MISSING_USER_ID` | 400 | X-User-ID header missing |
| `INSUFFICIENT_BALANCE` | 400 | Not enough funds |
| `WALLET_INACTIVE` | 400 | Wallet is not active |
| `INTERNAL_ERROR` | 500 | Server error |

---

## Rate Limits

- **Authentication:** 10 requests per minute
- **Wallet Operations:** 100 requests per minute
- **Transaction Operations:** 100 requests per minute
- **Webhook Operations:** 20 requests per minute

Rate limit headers:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1645123456
```

---

## Best Practices

### 1. **Use Idempotency Keys**
Always include `Idempotency-Key` header for transaction requests to prevent duplicate charges during retries.

```bash
Idempotency-Key: $(uuidgen)
```

### 2. **Implement Webhook Verification**
Always verify webhook signatures to ensure requests are from sample-provider.

### 3. **Handle Async Processing**
Transactions are asynchronous:
1. Create transaction (returns `PENDING`)
2. Poll transaction status OR wait for webhook
3. Handle `SUCCESS` or `FAILED` status

### 4. **Token Refresh**
Tokens expire in 1 hour. Implement token refresh logic:
```javascript
if (response.status === 401) {
  token = await refreshToken()
  retry(request)
}
```

### 5. **Error Handling**
Implement proper retry logic with exponential backoff for transient errors (500, 503).

### 6. **Logging**
Log all API interactions with:
- Request ID (from `X-Request-ID` response header)
- Transaction ID
- Timestamps
- Error details

---

## SDK Examples

### JavaScript/Node.js

```javascript
const axios = require('axios');

class sample-providerClient {
  constructor(clientId, clientSecret) {
    this.baseURL = 'https://api.sample-provider.co/v1';
    this.clientId = clientId;
    this.clientSecret = clientSecret;
    this.token = null;
  }

  async authenticate() {
    const response = await axios.post(`${this.baseURL}/auth/token`, {
      clientId: this.clientId,
      clientSecret: this.clientSecret
    });
    this.token = response.data.accessToken;
    return this.token;
  }

  async buyCredit(userId, walletId, amount, currency, metadata) {
    const response = await axios.post(
      `${this.baseURL}/transactions/credit-purchase`,
      { walletId, amount, currency, metadata },
      {
        headers: {
          'Authorization': `Bearer ${this.token}`,
          'X-User-ID': userId,
          'Idempotency-Key': `${Date.now()}-${Math.random()}`
        }
      }
    );
    return response.data.transaction;
  }

  async getTransactionStatus(transactionId) {
    const response = await axios.get(
      `${this.baseURL}/transactions/${transactionId}`,
      {
        headers: { 'Authorization': `Bearer ${this.token}` }
      }
    );
    return response.data.transaction;
  }
}

// Usage
const client = new sample-providerClient('bella_mobile_prod', 'secret_xyz');
await client.authenticate();

const transaction = await client.buyCredit(
  'usr_123',
  'wlt_usd_abc123',
  '10.00',
  'USD',
  { phoneNumber: '+1234567890' }
);

console.log('Transaction ID:', transaction.id);
```

### Python

```python
import requests
from typing import Dict, Optional

class sample-providerClient:
    def __init__(self, client_id: str, client_secret: str):
        self.base_url = 'https://api.sample-provider.co/v1'
        self.client_id = client_id
        self.client_secret = client_secret
        self.token: Optional[str] = None

    def authenticate(self) -> str:
        response = requests.post(
            f'{self.base_url}/auth/token',
            json={
                'clientId': self.client_id,
                'clientSecret': self.client_secret
            }
        )
        response.raise_for_status()
        self.token = response.json()['accessToken']
        return self.token

    def buy_credit(
        self,
        user_id: str,
        wallet_id: str,
        amount: str,
        currency: str,
        metadata: Dict[str, str],
        idempotency_key: Optional[str] = None
    ) -> Dict:
        headers = {
            'Authorization': f'Bearer {self.token}',
            'X-User-ID': user_id
        }
        if idempotency_key:
            headers['Idempotency-Key'] = idempotency_key

        response = requests.post(
            f'{self.base_url}/transactions/credit-purchase',
            json={
                'walletId': wallet_id,
                'amount': amount,
                'currency': currency,
                'metadata': metadata
            },
            headers=headers
        )
        response.raise_for_status()
        return response.json()['transaction']

# Usage
client = sample-providerClient('bella_mobile_prod', 'secret_xyz')
client.authenticate()

transaction = client.buy_credit(
    user_id='usr_123',
    wallet_id='wlt_usd_abc123',
    amount='10.00',
    currency='USD',
    metadata={'phoneNumber': '+1234567890'}
)

print(f"Transaction ID: {transaction['id']}")
```

---

## Support

For API support and issues:
- Email: api-support@sample-provider.co
- Developer Portal: https://developers.sample-provider.co
- Status Page: https://status.sample-provider.co
