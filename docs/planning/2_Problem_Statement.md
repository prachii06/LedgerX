
## Background

In distributed systems, a single business transaction is often processed by multiple independent services.

For example:

Customer places an order.

↓

Order Service creates the order.

↓

Payment Service confirms payment.

↓

Accounting Service records the transaction.

Since these services operate independently, events may not always arrive as expected.

## Problems

The system must handle situations where events:

- arrive late
- arrive more than once
- arrive in the wrong order
- never arrive

Without an automated reconciliation process, these situations can lead to:

- incorrect financial records
- duplicate transactions
- failed reconciliation
- manual investigation

## Proposed Solution

LedgerX continuously processes incoming events, detects inconsistencies, and reconciles related events to maintain a consistent transaction state.