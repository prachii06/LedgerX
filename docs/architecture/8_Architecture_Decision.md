# Architecture Decisions

## ADR-01: Go for Backend

Reason:
- Excellent concurrency support
- Simple deployment
- Strong performance
- Large ecosystem

---

## ADR-02: Kafka for Messaging

Reason:
- Durable event storage
- Reliable event streaming
- Decouples producers and consumers

---

## ADR-03: PostgreSQL as Source of Truth

Reason:
- ACID compliance
- Reliable persistence
- Strong relational capabilities

---

## ADR-04: Redis for Temporary State

Reason:
- Fast lookups
- Efficient handling of out-of-order events
- Reduces database load

---

## ADR-05: Docker Compose for Deployment

Reason:
- Simple setup
- Easy local development
- Suitable for MVP deployment on a VPS