## Overview

The database consists of four primary entities.

```text
+------------------+
|   Transactions   |
+------------------+
| transaction_id PK|
| status           |
| created_at       |
| updated_at       |
+---------+--------+
          |
          | 1
          |
          | *
+---------v--------+
|      Events      |
+------------------+
| event_id PK      |
| transaction_id FK|
| event_type       |
| source_service   |
| event_time       |
| processed        |
+------------------+

          |

          |

+----------------------+
| Reconciliation_Record|
+----------------------+
| transaction_id PK/FK |
| status               |
| missing_events       |
| last_updated         |
+----------------------+

          |

          |

+------------------+
|  Failed_Events   |
+------------------+
| id PK            |
| event_id         |
| reason           |
| retry_count      |
| created_at       |
+------------------+
```