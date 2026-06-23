# Overview

LedgerX is a transaction reconciliation system.

Its primary responsibility is to receive transaction events from multiple independent services, determine whether all required events for a transaction have been received, and maintain the correct transaction state.

The system operates using an event-driven architecture where events are processed asynchronously.

Instead of processing complete transactions directly, LedgerX processes individual events and gradually builds the final transaction state.

---

## Core Domain Objects

The main business objects are:

- Transaction
- Event
- Reconciliation Process
- Worker
- Dead Letter Queue (DLQ)

Each object has a specific responsibility within the reconciliation process.

---

## Business Goal

The business goal is to ensure that every transaction reaches a valid final state while handling:

- duplicate events
- delayed events
- out-of-order events
- missing events
- processing failures