# Error Handling

Luna SDKs use a consistent error hierarchy to help you handle failures predictably.

## Error Hierarchy

All errors thrown by the SDK inherit from a base `LunaError`.

| Error Type | Description | Retroable? |
|------------|-------------|------------|
| `AuthenticationError` | Invalid API key or expired token. | ❌ |
| `AuthorizationError` | Valid credentials but insufficient permissions. | ❌ |
| `ValidationError` | Invalid parameters provided. | ❌ |
| `NotFoundError` | The requested resource does not exist. | ❌ |
| `RateLimitError` | Too many requests. SDK automatically handles retries if configured. | ✅ |
| `NetworkError` | Connection failures or timeouts. | ✅ |
| `ServerError` | 5xx errors from the API side. | ✅ |

## Handling Errors

### TypeScript

```typescript
import { AuthenticationError, NotFoundError } from '@eclipse/luna-sdk/errors';

try {
  await client.users.get('usr_123');
} catch (err) {
  if (err instanceof NotFoundError) {
    console.log('User not found');
  } else if (err instanceof AuthenticationError) {
    console.log('Check your API key');
  } else {
    console.error('Unexpected error:', err);
  }
}
```

### Python

```python
from luna.errors import NotFoundError, AuthenticationError

try:
    client.users.get("usr_123")
except NotFoundError:
    print("User not found")
except AuthenticationError:
    print("Check your API key")
```

### Go

```go
import (
    "errors"
    lunaerr "github.com/eclipse-softworks/luna-sdk-go/luna/errors"
)

_, err := client.Users().Get(ctx, "usr_123")
var notFoundErr *lunaerr.NotFoundError
if errors.As(err, &notFoundErr) {
    fmt.Println("User not found")
}
```
