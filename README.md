# StockDB - Financial Market Database

## Design Overview
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
