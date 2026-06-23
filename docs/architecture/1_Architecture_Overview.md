## Overview

LedgerX follows an event-driven microservice-inspired architecture.

The system is composed of independent components that communicate through events rather than direct calls. This improves scalability, fault tolerance, and loose coupling.

The application consists of:

- API Service
- Kafka
- Reconciliation Worker
- PostgreSQL
- Redis
- Dead Letter Queue (DLQ)
- Operations Dashboard

Each component has a single responsibility and can be developed, tested, and deployed independently.

---

## Architecture Style

LedgerX uses the following architectural principles:

- Event-Driven Architecture
- Producer–Consumer Pattern
- Asynchronous Processing
- Stateless API Service
- Idempotent Event Processing
- Layered Backend Design

---

## Benefits

- Loose coupling between services
- Better scalability
- Improved fault tolerance
- Easier horizontal scaling
- Clear separation of responsibilities