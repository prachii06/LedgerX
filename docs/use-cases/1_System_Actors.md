## 1. Transaction Generator

Simulates real-world systems that produce transaction events.

Examples:
- Order Service
- Payment Gateway
- Accounting Service

Role:
- Sends transaction-related events to the system.

---

## 2. LedgerX API Service

The entry point of the system.

Role:
- Receives events
- Validates events
- Forwards events to Kafka

---

## 3. Kafka (Event Broker)

Role:
- Stores and streams events
- Decouples producers and consumers

---

## 4. Reconciliation Worker

Background service that processes events.

Role:
- Consumes events from Kafka
- Applies reconciliation logic
- Updates Redis and PostgreSQL

---

## 5. Redis

Role:
- Stores temporary reconciliation state
- Helps match delayed/out-of-order events quickly

---

## 6. PostgreSQL

Role:
- Stores final transaction state
- Acts as system of record

---

## 7. Operations Dashboard (Frontend)

Role:
- Displays system state
- Shows transactions, events, and metrics
- Helps monitor system health