## Performance

- Event processing should be asynchronous.
- API responses should be fast.
- Redis should reduce reconciliation lookup time.

---

## Scalability

- The system should support multiple reconciliation workers.
- Components should be independently scalable in the future.

---

## Reliability

- Events should not be lost during processing.
- Duplicate processing should be prevented.

---

## Availability

- The system should continue processing even if individual events fail.
- Failed events should be isolated using a Dead Letter Queue.

---

## Maintainability

- The codebase should follow a modular architecture.
- Components should be loosely coupled.

---

## Observability

The system should expose metrics for monitoring.

Metrics include:

- Event throughput
- Processing latency
- Failure count
- Duplicate count
- Pending transactions

---

## Security

- Input validation
- Environment-based configuration
- Secure handling of secrets

---

## Deployability

The complete application should be deployable using Docker Compose on a single VPS.