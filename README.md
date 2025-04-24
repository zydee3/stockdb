# StockDB - Personal Financial Market Database

StockDB is a lightweight financial market data collection and storage solution
designed for personal or small-scale environments. The system follows an
intuitive data model that mirrors how individuals naturally think about stock
information, ensuring queries are both intuitive and efficient.

## StockDB Architecture

StockDB follows a Kubernetes-inspired architecture with two primary components:

1. `stockd`: A daemon service that orchestrates data collection through a
   manager and worker pool
2. `stockctl`: A command-line interface tool for interacting with the daemon

Users submit data collection jobs by creating YAML configuration files and
applying them through stockctl, similar to how Kubernetes handles Custom
Resource Definitions (CRDs).

### Stockd Service

1. Manager:
   - Processes and validates incoming job definitions
   - Schedules jobs based on priority and resource availability
   - Monitors job execution and handles failure scenarios
   - Maintains job state and provides status updates
   - Implements configurable retry policies for failed jobs
   - Provides persistence to protect against service interruptions
2. Job Queue: A light weight priority queue that:
   - Maintains prioritized pending and in-progress jobs
   - Supports configurable job priorities
   - Supports efficient job batching for optimal resource utilization
3. Workers: A job processor that:
   - Executes data collection jobs from multiple financial data sources
   - Standardizes collected data into a consistent internal format
   - Implements rate limiting and backoff strategies for API calls
   - Performs data validation before storage
   - Optimizes batch writes to the database
4. Communication
   - Listenes on a Unix socket for local communication with `stockctl`
   - Implements well-defined communication for job submission and status queries
   - Provides graceful shutdown and request metrics for monitoring and logging

### Stockctl Command-line Tool

- Creates, submits, monitors, and manages data collection jobs
- Validates job configurations before submission
- Supports real-time status monitoring of running jobs
- Offers job management capabilities (pause, resume, cancel)

## Workflow Examples

### Successful Job Processing

1. The user submits a CRD to collect data for a security.
2. The manager creates separate job batches for each security.
3. The worker claims a batch of jobs that are related by security.
4. The worker processes each job in the batch (collect and standardize).
5. The worker writes the standardized data batch to the database.
6. The worker marks the job as completed and provides a completion certificate.

### Handling Failures.

1. If a job fails (API unavailable, network issues, etc), it will be retried.
2. The worker writes the standardized data from the successful jobs to the
   database.
3. The worker marks the job as incomplete and provides a completion certificate
   for the successfully completed jobs.
4. The manager reclaims the failed job.
5. Failed job enters an exponential backoff period and retried.
6. After a configurable retry limit, the job is logged and abandoned.

## Database Design Overview

StockDB implements a securities-anchored architecture using TimescaleDB
hypertables with time-based partitioning. There were two approaches considered
when designing StockDB:

1. **Date-Anchored Architecture**

   - Uses trading days as central reference point
   - Organizes all data around specific calendar days
   - Simplifies calendar-based operations and data association

2. **Securities-Anchored Architecture** (Selected)
   - Uses securities as the primary dimension
   - Employs TimescaleDB hypertables with time as partition key
   - Maintains natural relationship between financial entities

We selected the securities-anchored approach for:

- Natural alignment with financial domain relationships. Users, such as market
  analysts tend to think in terms of "AAPL performance" or "NVDA earnings" rather
  than "January 15th market events."
- Efficient support for securities-centric queries. When an analyst needs to
  compare Microsoft's price movements against their earnings surprises over the
  past 8 quarters, the query can leverage direct indexing rather than joining
  across multiple date-based tables.
- Better handling of irregular data points and preserves the continuity of
  market narratives that naturally span across multiple time periods. When Meta
  faces a regulatory investigation that unfolds over weeks, the
  securities-anchored model maintains the coherent storyline by linking all
  related news to META rather than fragmenting it across disconnected daily
  records. This approach recognizes that financial events rarely conform to neat
  daily boundaries and often require analysis across irregular time intervals.
- Better scaling with growing data volume. A securities-anchored architecture
  creates natural sharding boundaries by security, then time, allowing independent
  scaling and retention policies for different entities. Apple stock might
  generate 10x the trading volume and news coverage of a mid-cap stock, requiring
  different optimization strategies.

## Data Structure

- Securities (dimension table): Central anchor reference for all financial
  entities
- Fact tables (hypertable): Independent data points such as prices, news
  articles, earnings, etc.
- Junction tables: For many-to-many relationships (e.g., news-securities)
- External storage: For large unstructured content (e.g., full articles)

## Performance Optimization

- Composite indexes on (security_id, time)
- Conflict handling for parallel ingestion workers
- Continuous aggregates for common calculations
- Appropriate chunk sizing for query patterns
- Compression policies for older data

## Implementation Notes

- All fact tables use time as partition key
- ON CONFLICT clauses for upsert operations
- Hash-based existence checks for rapid duplicate detection
- Appropriate constraints to maintain data integrity
