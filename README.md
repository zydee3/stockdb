# StockDB - Personal Financial Market Database

StockDB is a lightweight financial market data collection and storage solution
designed for individual analysts and researchers. Data is stored in an intuitive
structure that aligns with how individuals naturally think about stock
information.

## Core Design Philosophy

StockDB prioritizes reliability and intuitive data organization within modest
resource constraints. Rather than attempting to compete with enterprise-scale
solutions, we focus on:

1. Simplicity: Clear, understandable architecture that's easy to maintain
2. Reliability: Simple and effective consistency and error handling.
3. Intuitive Design: The implementation should be intuitive.
4. Efficient: Optimized for personal or small-scale environments.

## System Architecture

StockDB implements a CRD-based scheduler (Kubernetes-inspired) architecture with
3 core components:

### 1. Manager

The manager handles job scheduling and processing:

- Accepts Custom Resource Definitions (CRDs) that define data collection tasks.
- Monitors job status and handles job completion failures.

### 2. Job Queue

A light weight priority queue that:

- Orders jobs by priority.
- Holds pending and incomplete jobs for workers to consume.

### 3. Workers

A job processor that:

- Accepts job batches from the job queue.
- Standardizes collected data into a consistent format.
- Writes completed batches to the database.
- Reports the job status to the Manager.

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
