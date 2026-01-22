# Mock Client

The `MockClient` allows you to develop against the Luna SDK API without needing a running backend.

## Usage

```typescript
import { MockClient, mockUser } from '@eclipse-softworks/luna-sdk/testing';

const client = new MockClient({
  // Seed initial data
  data: {
    users: {
      'usr_1': mockUser({ name: 'Alice' })
    }
  },
  // Simulate network conditions
  latency: 500,
  errorRate: 0.1
});

// Use just like the real client
const users = await client.users.list();
```
