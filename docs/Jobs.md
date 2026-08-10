# JOB permissions

| Action        | Receptionist  |     Technician     |  Head Technician   | Admin  | Super Admin  |
| ------------- | :----------:  | :---------------:  | :---------------:  | :---:  | :---------:  |
| Create Job    |       ✅      |         ❌         |         ✅         |   ✅   |      ❌      |
| Assign Job    |       ❌      |         ❌         |         ✅         |   ✅   |      ❌      |
| Change Status |       ❌      | ✅ (assigned only) | ✅ (assigned only) |   ❌   |      ❌      |
| Add Notes     |       ✅      | ✅ (assigned only) | ✅ (assigned only) |   ❌   |      ❌      |
| Close Job     |       ❌      | ✅ (assigned only) | ✅ (assigned only) |   ❌   |      ❌      |
| View Jobs     |       ✅      |         ✅         |         ✅         |   ✅   |      ✅      |


# Job Closure

| Closure reason               | Allowed status                                           |
| ---------------------------- | -------------------------------------------------------- |
| **Completed**                | `resumed`                                                |
| **Not repairable**           | `waiting_customer`, `resumed`                            |
| **Customer no response**     | `waiting_customer`, `resumed`                            |
| **Customer declined repair** | `waiting_customer`, `resumed`                            |
| **Customer cancelled**       | `in_progress`, `waiting_customer`, `resumed`             |
| **Duplicate job**            | `assigned`, `in_progress`, `waiting_customer`, `resumed` |


## Closure flowchart

created
   │
   │ assign
   ▼
assigned
   │
   └── duplicate_job ─────────→ closed
   │
   │ technician starts work
   ▼
in_progress
   │
   ├── customer_cancelled ───→ closed
   ├── duplicate_job ────────→ closed
   │
   │ waiting for customer
   ▼
waiting_customer
   │
   ├── not_repairable ───────→ closed
   ├── customer_no_response ─→ closed
   ├── customer_declined ────→ closed
   ├── customer_cancelled ───→ closed
   └── duplicate_job ────────→ closed
   │
   │ customer responds
   ▼
resumed
   │
   ├── completed ────────────→ closed
   ├── not_repairable ───────→ closed
   ├── customer_no_response ─→ closed
   ├── customer_declined ────→ closed
   ├── customer_cancelled ───→ closed
   └── duplicate_job ────────→ closed