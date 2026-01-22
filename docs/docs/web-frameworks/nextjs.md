# Next.js Integration

Luna SDK provides specialized adapters for Next.js 14+ with support for Server Components, API Routes, and Edge Middleware.

## Server Components

Use `createServerClient` in `page.tsx` or `layout.tsx`. Does not expose API keys to the client.

```tsx
import { createServerClient } from '@eclipse-softworks/luna-sdk/adapters';

export default async function UsersPage() {
  const client = createServerClient(); // Uses process.env.LUNA_API_KEY
  const users = await client.users.list({ next: { revalidate: 60 } });

  return (
    <ul>
      {users.items.map(user => (
        <li key={user.id}>{user.name}</li>
      ))}
    </ul>
  );
}
```

## API Routes

Create type-safe API handlers in `app/api/[...luna]/route.ts`.

```typescript
import { createApiHandler } from '@eclipse-softworks/luna-sdk/adapters';

export const { GET, POST, PATCH, DELETE } = createApiHandler({
    // Optional: add custom middleware
    middleware: [authMiddleware]
});
```

## Middleware

Use `createMiddleware` in `middleware.ts` for edge-compatible logic.

```typescript
import { createMiddleware } from '@eclipse-softworks/luna-sdk/adapters';
import { NextResponse } from 'next/server';

export const middleware = createMiddleware(async (req, client) => {
    // Verify a webhook signature or check a user permission at the edge
    return NextResponse.next();
});

export const config = {
    matcher: ['/api/:path*'],
};
```
