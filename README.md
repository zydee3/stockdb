# StockDB - Personal Financial Market Database

StockDB is a lightweight financial market data ingestion program that collects, 
standardizes, and stores market information. Data is stored in an intuitive 
structure that aligns with how individuals naturally think about stock 
information. While optimized for performance, StockDB makes strategic trade-offs 
to prioritize reliability within modest resource constraints. We deliberately 
favor simplicity and maintainability over complex enterprise-scale capabilities, 
creating a solution ideal for personalized self-hosted financial analytics.

## Implementation Overview

### Architecture Core Principles
1. Each component should be modular.
2. Each step in the pipeline should be durable (i.e., all failures must be 
recoverable).
3. Each component must be stateless.
4. Each component must be able to prove it's work was completed.

### Architecture Components
StockDB uses modular components that can run together or independently, allowing 
for easy replacement or horizontal scaling of individual modules should you need 
to scale. 

1. **Workers**: Source-specific collectors that connect external data sources to 
StockDB. Each worker fetches raw data from designated APIs, standardizes it, and 
pushes it to the Stream. Workers treat data sources as "black boxes" - focusing 
only on retrieval and standardization. The data fetched can be configured to 
persisted on disk until the Stream provides a valid certificate confirming the 
data has been successfully processed. 
2. **Stream**: Lightweight pub/sub streamer that implements data persistence to 
ensure durability. The data recieved can be configured to be stored on the disk
until the Stream receives a valid certificate confirming the data has been
received by all subscribers.
4. **Batch Processor**: Consumes standardized data from the stream and assembles 
efficient batches for database operations.
5. **Caches**: Caches are used to store recent query results. For example, a 
cache can be used to store the most recent queries handled by the worker to 
avoid redundant query processing. 

### Sample Workflow
1. Worker fetches raw data from an API.
2. Worker writes raw data to disk for persistence.
3. Worker standardizes the raw data into a standardized format.
4. Worker publishes the standardized data to the Stream.
5. Stream receives the standardized data and writes it to disk for persistence.
6. Stream writes the standardized data to disk for persistence.
7. Stream writes a certificate to Worker.
8. Worker receives the certificate and deletes the raw data from disk.
9. Stream writes the data to the topic subcribers.
10. Each subscriber receives the data, processes it, writes a certificate to the
Stream. 
11. Stream receives the certificates and deletes the data from disk.

### Sample Workflow Failure Recovery
1. If the system crashes during steps 3-8, the system can recovery by having the 
workers re-read the data from disk, re-standardizing, and re-publishing it to 
the Stream.
2. If the system crashes during steps 9-11, the system can recovery by having
the Stream re-read the data from disk and re-publish it to the subscribers.

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
