## Transaction States

NEW

↓

PENDING

↓

PROCESSING

↓

RECONCILING

↓

RECONCILED

OR

FAILED

OR

PENDING_REVIEW

---

## State Descriptions

### NEW

Transaction has been created but no events have been processed.

---

### PENDING

Waiting for required events.

---

### PROCESSING

Worker is currently processing incoming events.

---

### RECONCILING

System is verifying whether all expected events have been received.

---

### RECONCILED

All required events were received successfully.

Final State.

---

### FAILED

Transaction cannot be reconciled due to unrecoverable errors.

Final State.

---

### PENDING_REVIEW

Transaction requires manual investigation because expected events were not received within the timeout period.

Final State.