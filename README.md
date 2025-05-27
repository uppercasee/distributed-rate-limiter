# Distributed Rate Limiting Service (DRLS)

Rate limiting is a common mechanism used in modern systems to control the rate at which clients can access resources. It protects services from abuse, ensures fair usage, and helps manage traffic during high-load scenarios.

The Distributed Rate Limiting Service (DRLS) is designed to provide scalable, consistent rate limiting across multiple service instances. It is implemented as a gRPC service and backed by Redis, using a sliding window algorithm to track request activity in near real-time.

## Core Concept

At its core, DRLS enforces limits on how many requests a specific client can make within a given time window. Instead of tracking counters per fixed interval, DRLS uses a **sliding window** approach, where each request is timestamped and stored temporarily in Redis. This allows it to provide smoother enforcement, avoiding burstiness allowed by fixed windows.

For example, if a client is allowed 100 requests per minute, the system ensures that the client does not exceed 100 requests in any 60-second rolling window, regardless of when the minute starts.

## Redis as Central Store

Redis is used as the single source of truth for request timestamps. All gRPC server replicas connect to the same Redis instance, which stores request data using sorted sets (`ZADD`, `ZREMRANGEBYSCORE`). This enables multiple server replicas to consistently enforce limits while sharing state across the cluster.

Each request results in a timestamp being added to Redis under a key representing the client. Old timestamps are pruned, and the current request count is compared to the limit.

## Distributed gRPC Servers

To scale the service horizontally, DRLS supports multiple gRPC server instances running in parallel. These are fronted by a load balancer (e.g., NGINX) that distributes traffic. All replicas are stateless and identical — they rely entirely on Redis to enforce limits, making it easy to scale up or down depending on load.

This architecture ensures that the rate limit is globally enforced, regardless of which replica handles the request, as long as they share access to the same Redis instance.

## Sliding Window Algorithm

The sliding window algorithm used by DRLS avoids the abrupt resets that can occur with fixed windows. Instead of counting requests per strict interval (e.g., 12:00–12:01), it looks back over the last N seconds from the time of each incoming request and counts how many requests occurred during that period.

This provides more accurate enforcement and fairness, and helps smooth out spikes in traffic that could otherwise be allowed under a fixed window model.

## TODO / Future Improvements

- ~~Use **Redis Lua scripts** for atomic request handling to prevent race conditions.~~
- Add support for **Redis Cluster / sharding** for better scalability and fault tolerance.
- Integrate **Prometheus metrics exporter** for observability and performance monitoring.
- Add **configurable rate limit policies** per client or API key.
- Support **graceful Redis failover** and connection retry logic for high availability.
