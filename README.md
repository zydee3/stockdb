# StockDB - Personal Financial Market Database

StockDB is a purpose-built, lightweight database solution for standardizing and 
storing financial market data with an intuitive organization that aligns with 
how analysts naturally think about financial information. StockDB maintains a 
high-performance architecture while making deliberate trade-offs that favor 
simplicity and maintainability over enterprise-scale capabilities. Our approach 
aims to deliver reliable performance within modest resource constraints by using 
lightweight, integrated solutions rather than complex external systems meant for
enterprise-level scaling.

## Implementation Overview

### Architecture Components
StockDB uses modular components that can run together or independently, allowing 
for easy replacement or horizontal scaling of individual modules should you need 
to scale.

1. **Workers**: Source-specific data collectors that fetch from various APIs, 
standardize the data into a common format, and push it to the queue.
2. **Queue**: Lightweight Go channel-based queue with disk persistence that 
buffers standardized data.
3. **Persistence Layer**: Implements write-ahead logging to ensure data 
durability before database operations.
4. **Batch Processor**: Consumes standardized data from the queue and assembles 
efficient batches for database operations.
5. **Simple Caches**: Simple caches can be used at the worker and database 
level. A worker cache can help coordinate redundant decision paths or data point processing. A database cache can be used to store recent query results. The size 
of each cache is configurable.

### Key Design Decisions
1. Each component is designed to be modular and replacable if needed. 
2. Data is persisted to prevent loss in case of failure at any stage in the 
pipeline until the data is confirmed to be in the database. 
3. Batching is triggered by either time or batch size. 
4. Failed operations can be retrieved and retried from persistence storage. 

### Sample Workflow
1. A scraper worker fetches the latest price data for AAPL.
2. Worker standardizes the raw price data for AAPL.
3. Worker pushes the standardized data to the Queue.
4. Persistence Layer writes the standardized data to disk (write-ahead log).
5. Batch Processor pulls items from the Queue, accumlating items until reaching
the configured threshold (e.g., 10 seconds or 100 items).
6. Batch Processor writes the batch to the database and generates a 
confirmation.
7. Persistence accepts the confirmation and cleans up the processed items.

If the system crashes during step 6, the system can recovery due to the 
persistence layer. Items not written to the database will be stored on the disk
and reconsumed by the Batch Processor upon system start.

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
