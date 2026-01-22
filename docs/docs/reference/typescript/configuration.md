# Configuration

## LunaConfig interface

```typescript
interface LunaConfig {
  apiKey: string;
  baseUrl?: string;
  timeout?: number;
  retries?: number;
  cache?: CacheConfig;
  metrics?: MetricsConfig;
}
```
