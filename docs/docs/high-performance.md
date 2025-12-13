---
sidebar_position: 6
---

# High Performance

The Luna SDK v2.0 introduces a high-performance architecture inspired by modern framework best practices ("Spring Boot/Rust/C").

## 1. Auto-Configuration ("Spring Boot")

The SDK automatically detects configuration from environment variables, allowing for zero-code initialization.

**Supported Variables:**
*   `LUNA_API_KEY`
*   `LUNA_ACCESS_TOKEN`
*   `LUNA_BASE_URL`
*   `LUNA_ENV`

### Example
**Environment**:
```bash
export LUNA_API_KEY=lk_live_xxxx
```

**Code**:
```python
# No arguments needed!
client = LunaClient()
```

## 2. Strict Mode ("Rust")

Enable **Strict Mode** to enforce rigorous client-side validation. This prevents invalid data from ever reaching the network, saving latency and API costs.

*   **Checks performed**: ID number checksums, Tax reference formats, Registration number regex.

### Enabling Strict Mode

**Go**:
```go
client, _ := luna.NewClient(luna.WithStrictMode())
```

**Python**:
```python
client = LunaClient(strict=True)
```

## 3. Connection Pooling ("C")

The SDKs are pre-configured with optimized HTTP transport settings for high-throughput applications.

*   **Go**: `MaxIdleConns: 100`, `IdleConnTimeout: 90s`
*   **Python**: `httpx` limits tuned to `max_keepalive=20`, `max_connections=100`

No manual configuration is required to benefit from these optimizations.
