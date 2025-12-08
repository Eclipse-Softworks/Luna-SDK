# @eclipse/luna-sdk

Official TypeScript SDK for the Eclipse Softworks Platform API.

## Installation

```bash
npm install @eclipse/luna-sdk
```

## Quick Start

```typescript
import { LunaClient } from '@eclipse/luna-sdk';

// API Key authentication
const client = new LunaClient({
  apiKey: process.env.LUNA_API_KEY,
});

// List users
const users = await client.users.list({ limit: 10 });

// Get a specific user
const user = await client.users.get('usr_123');

// Create a new user
const newUser = await client.users.create({
  email: 'john@example.com',
  name: 'John Doe',
});
```

## Authentication

### API Key

```typescript
const client = new LunaClient({
  apiKey: 'lk_live_xxxx',
});
```

### OAuth Token

```typescript
const client = new LunaClient({
  accessToken: session.accessToken,
  refreshToken: session.refreshToken,
  onTokenRefresh: async (tokens) => {
    await saveTokens(tokens);
  },
});
```

## Error Handling

```typescript
import { LunaClient, NotFoundError, RateLimitError } from '@eclipse/luna-sdk';

try {
  await client.users.get('usr_nonexistent');
} catch (error) {
  if (error instanceof NotFoundError) {
    console.log('User not found:', error.message);
  } else if (error instanceof RateLimitError) {
    console.log('Rate limited, retry after:', error.retryAfter);
  }
}
```

## Configuration

```typescript
const client = new LunaClient({
  apiKey: process.env.LUNA_API_KEY,
  baseUrl: 'https://api.staging.eclipse.dev', // Custom base URL
  timeout: 60000, // Request timeout (ms)
  maxRetries: 5, // Retry attempts
  logLevel: 'debug', // Log level
});
```

## License

MIT © Eclipse Softworks
