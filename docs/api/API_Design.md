# API Design

## Overview

The LedgerX API provides REST endpoints for external systems and the Operations Dashboard to interact with the platform.

The API is responsible for:

- Receiving transaction events
- Retrieving transaction details
- Viewing reconciliation status
- Accessing system metrics
- Monitoring application health

The API communicates using JSON over HTTP.

---

# Base URL

```
/api/v1
```

---

# Authentication

Authentication is not included in Version 1.

Future versions may support:

- JWT Authentication
- API Keys
- Role-Based Access Control (RBAC)

---

# Content Type

Request

```
application/json
```

Response

```
application/json
```

---

# API Conventions

## Resource Naming

Use plural nouns.

Examples:

```
/events
/transactions
/dashboard
```

---

## HTTP Methods

| Method | Purpose |
|---------|----------|
| GET | Retrieve resources |
| POST | Create resources |
| PUT | Replace resources |
| PATCH | Partial update |
| DELETE | Remove resources |

---

## HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | Success |
| 201 | Resource Created |
| 400 | Bad Request |
| 404 | Not Found |
| 409 | Conflict (Duplicate Event) |
| 500 | Internal Server Error |

---

# Event APIs

## POST /events

Receives a transaction event.

### Request

```json
{
  "event_id": "uuid",
  "transaction_id": "uuid",
  "event_type": "PAYMENT_COMPLETED",
  "source_service": "payment-service",
  "timestamp": "2026-06-23T10:30:00Z"
}
```

### Success Response

```json
{
  "message": "Event accepted",
  "status": "queued"
}
```

---

## GET /events/{eventId}

Returns details of a specific event.

### Response

```json
{
  "event_id": "uuid",
  "transaction_id": "uuid",
  "event_type": "PAYMENT_COMPLETED",
  "processed": true
}
```

---

# Transaction APIs

## GET /transactions

Returns a paginated list of transactions.

Example

```
GET /transactions?page=1&limit=20
```

---

## GET /transactions/{transactionId}

Returns details of a transaction.

### Response

```json
{
  "transaction_id": "uuid",
  "status": "RECONCILED",
  "created_at": "...",
  "updated_at": "..."
}
```

---

## GET /transactions/{transactionId}/events

Returns all events associated with a transaction.

### Response

```json
[
  {
    "event_id": "...",
    "event_type": "ORDER_CREATED"
  },
  {
    "event_id": "...",
    "event_type": "PAYMENT_COMPLETED"
  }
]
```

---

# Dashboard APIs

## GET /dashboard/summary

Returns overall system statistics.

### Response

```json
{
  "transactions": 1250,
  "reconciled": 1180,
  "pending": 45,
  "failed": 25
}
```

---

## GET /dashboard/metrics

Returns operational metrics.

### Response

```json
{
  "events_processed": 25000,
  "processing_rate": 120,
  "worker_status": "healthy"
}
```

---

# Health APIs

## GET /health

Returns application health.

### Response

```json
{
  "status": "healthy"
}
```

---

## GET /ready

Returns application readiness.

### Response

```json
{
  "status": "ready"
}
```

---

## GET /live

Returns liveness status.

### Response

```json
{
  "status": "alive"
}
```

---

# Standard Error Response

All errors follow the same structure.

```json
{
  "error": {
    "code": "INVALID_EVENT",
    "message": "Event type is invalid"
  }
}
```

---

# Common Error Codes

| Code | Description |
|------|-------------|
| INVALID_EVENT | Invalid event payload |
| DUPLICATE_EVENT | Event already processed |
| TRANSACTION_NOT_FOUND | Transaction does not exist |
| INTERNAL_ERROR | Unexpected server error |

---

# API Versioning

LedgerX uses URL-based versioning.

Examples

```
/api/v1/events
/api/v2/events
```

Breaking changes require a new API version.

---

# Future APIs

The following endpoints may be added in future versions:

- Authentication APIs
- Replay Failed Events
- Dead Letter Queue Management
- Event Replay
- Audit Logs
- Worker Management