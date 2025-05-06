# StockDB Manager Job Management

The manager accepts jobs and manages them with dynamic job prioritization and 
Kubernetes-inspired architecture.

### Components
- **stockd Daemon:**  
  Handles job scheduling, queue management, and worker orchestration.  

- **stockctl CLI:**  
  User-facing tool for submitting, querying, and managing jobs.

- **SQLite Storage:**  
  Embedded, persistent job store with optimized schema and indexes.


## Features
- **Temporal Job Queues:**  
  Real-time jobs (created today, UTC) and historical jobs (created before today) 
  are automatically classified and processed in optimal order.

- **Descending Real-Time Processing:**  
  Real-time jobs are always processed newest-first.

- **ACID-Compliant Embedded Storage:**  
  SQLite database with WAL mode for safe concurrent access.

- **Exponential Backoff & Retries:**  
  Configurable retry logic for failed jobs.

- **Kubernetes-Inspired Workflow:**  
  Job definitions via YAML, submitted and managed using a CLI.
  
- **Efficient Indexing:**  
  Partial indexes for fast queue operations.


## Database Schema
### jobs Table

| Column        | Type      | Description                               |
|---------------|-----------|-------------------------------------------|
| id            | INTEGER   | Primary key                               |
| job_id        | TEXT      | Unique job identifier (from YAML)         |
| job_type      | TEXT      | 'RECURRING' or 'INTERVAL'                 |
| status        | TEXT      | 'pending', 'processing', 'completed', 'failed' |
| spec_json     | TEXT      | Full job definition (as JSON)             |
| created_at    | DATETIME  | When job was created (UTC)                |
| next_run_time | DATETIME  | Next scheduled execution time             |
| last_updated  | DATETIME  | Last status update                        |
| attempts      | INTEGER   | Number of attempts                        |
| max_retries   | INTEGER   | Maximum allowed retries                   |
| schedule_params | TEXT    | JSON-encoded schedule details             |

---

#### Indexes
```sql
-- Real-time queue: pending jobs created today, newest first
CREATE INDEX idx_realtime ON jobs (created_at DESC)
WHERE status = 'pending' AND date(created_at) = date('now');

-- Historical queue: pending jobs created before today, oldest first
CREATE INDEX idx_historical ON jobs (created_at ASC)
WHERE status = 'pending' AND date(created_at) < date('now');
```


## How Real-Time and Historical Queues Work
- **Real-Time Queue:**  
  Jobs with `created_at` equal to today's UTC date.  
  Processed in descending order (`newest first`).

- **Historical Queue:**  
  Jobs with `created_at` before today (UTC).  
  Processed in ascending order (`oldest first`).

- **No manual priority field:**  
  Temporal classification and queue ordering are automatic.

## Example Job Yamls

There are two type of jobs: one time jobs and repeated jobs. 

### Collection Job
There are two 