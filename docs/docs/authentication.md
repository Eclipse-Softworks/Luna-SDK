# Authentication

The Luna SDK supports both **API Key** (machine-to-machine) and **OAuth Token** (user-centric) authentication.

## API Key Authentication

Best for server-side scripts, cron jobs, and backend integration.

### TypeScript
```typescript
import { LunaClient } from '@eclipse/luna-sdk';

const client = new LunaClient({
  apiKey: 'lk_live_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6'
});
```

### Python
```python
from luna import LunaClient

client = LunaClient(api_key="lk_live_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6")
```

### Go
```go
import "github.com/eclipse-softworks/luna-sdk-go/luna"

client := luna.NewClient(luna.WithAPIKey("lk_live_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"))
```

## OAuth Token Authentication

Best for client-side apps, CLIs, and applications acting on behalf of a user. The SDK handles token refresh automatically.

### TypeScript
```typescript
const client = new LunaClient({
  accessToken: 'ey...',
  refreshToken: 'rt...',
  // Optional: Callback when tokens are refreshed
  onTokenRefresh: (auth) => {
    localStorage.setItem('tokens', JSON.stringify(auth));
  }
});
```

### Python
```python
client = LunaClient(
    access_token="ey...",
    refresh_token="rt...",
    on_token_refresh=save_tokens_callback
)
```

### Go
```go
client := luna.NewClient(
    luna.WithTokens("ey...", "rt..."),
    luna.WithTokenRefreshCallback(func(tokens auth.TokenPair) error {
        // Save new tokens
        return nil
    }),
)
```

## Security Best Practices

1. **Never commit API keys** to version control. Use environment variables.
2. **Use specific API keys** with scoped permissions for different services.
3. **Rotate keys** immediately if you suspect a leak.
