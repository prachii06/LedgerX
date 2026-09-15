# LedgerX Load Testing

This directory contains k6 load testing scripts for the LedgerX project.

## Prerequisites
- [k6](https://k6.io/docs/get-started/installation/) must be installed.
- LedgerX backend, PostgreSQL, Kafka, and Redis must be running locally.

## Running Tests

To run the transactions load test (tests general HTTP API limits and async event generation):
```bash
k6 run .\load-tests\transactions.js
```

To run the simulation load test (tests bulk event generation):
```bash
k6 run .\load-tests\simulation.js
```

To run the reconciliation test (tests full end-to-end processing pipeline including database queries):
```bash
k6 run .\load-tests\reconciliation.js
```

## Options
You can override the API URL by passing `API_URL` environment variable:
```bash
k6 run -e API_URL=http://localhost:8080 .\load-tests\transactions.js
```
