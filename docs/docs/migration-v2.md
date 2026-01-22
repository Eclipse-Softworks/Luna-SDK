# Migration Guide: v1.x to v2.0

This guide helps you migrate from Luna SDK v1.x to v2.0.

## Overview

Luna SDK v2.0 introduces significant new features while maintaining backward compatibility for core functionality. Most v1.x code will work without changes.

## Breaking Changes

### Go Module Path

The Go module path now includes a `/v2` suffix for proper semantic versioning:

```diff
- import "github.com/eclipse-softworks/luna-sdk-go/luna"
+ import "github.com/eclipse-softworks/luna-sdk-go/v2/luna"
```

```diff
- go get github.com/eclipse-softworks/luna-sdk-go
+ go get github.com/eclipse-softworks/luna-sdk-go/v2
```

### Minimum Requirements

| Requirement | v1.x | v2.0 |
|-------------|------|------|
| Node.js | 16+ | 18+ |
| Python | 3.8+ | 3.9+ |
| Go | 1.20+ | 1.21+ |

---

## New Features in v2.0

### South African Payments

New payment gateway integrations:

```typescript
const client = new LunaClient({
    apiKey: process.env.LUNA_API_KEY,
    payments: {
        payfast: { merchantId: '...', merchantKey: '...' },
        ozow: { siteCode: '...', privateKey: '...' },
        yoco: { secretKey: '...' },
        payshap: { participantCode: '...' },
    },
});

// Create PayFast payment
const payment = await client.payments.payfast.create({
    amount: { value: 100, currency: 'ZAR' },
    reference: 'ORDER-123',
    returnUrl: 'https://example.com/success',
});
```

### Messaging (SMS, WhatsApp, USSD)

```typescript
const client = new LunaClient({
    apiKey: process.env.LUNA_API_KEY,
    messaging: {
        sms: { provider: 'clickatell', apiKey: '...' },
        whatsapp: { accessToken: '...', phoneNumberId: '...' },
    },
});

// Send SMS
await client.messaging.sms.send({
    to: '+27821234567',
    message: 'Hello from Luna!',
});
```

### AI/ML Capabilities

```typescript
// Text embeddings
const embeddings = await client.ai.embeddings.create({
    input: 'Hello world',
    model: 'text-embedding-ada-002',
});

// Vision analysis
const analysis = await client.ai.vision.analyze({
    imageUrl: 'https://example.com/image.jpg',
    prompt: 'Describe this image',
});

// SA language translation
const translated = await client.ai.translate({
    text: 'Hello, how are you?',
    targetLanguage: 'isiZulu',
});
```

### South African Business Tools

```typescript
// CIPC company lookup
const company = await client.zaTools.cipc.lookup('2020/123456/07');

// B-BBEE verification
const certificate = await client.zaTools.bbbee.verify('CERT-123');

// SA ID validation
const idInfo = client.zaTools.id.parse('8501015009087');
console.log(idInfo.dateOfBirth); // 1985-01-01
console.log(idInfo.gender); // 'male'
console.log(idInfo.isCitizen); // true
```

---

## New HTTP Utilities

### Circuit Breaker

Prevent cascading failures:

```typescript
import { CircuitBreaker, CircuitState } from '@eclipse-softworks/luna-sdk';

const breaker = new CircuitBreaker({
    failureThreshold: 5,
    resetTimeout: 30000,
    onStateChange: (from, to) => {
        console.log(`Circuit: ${from} -> ${to}`);
    },
});

if (!breaker.canRequest()) {
    throw new Error('Service unavailable');
}
```

### Rate Limiter

Handle rate limits gracefully:

```typescript
import { RateLimiter } from '@eclipse-softworks/luna-sdk';

const limiter = new RateLimiter({
    maxWaitTime: 60000,
    onRateLimited: (state, waitTime) => {
        console.log(`Rate limited, waiting ${waitTime}ms`);
    },
});

// After each request, update from response headers
limiter.updateFromHeaders(response.headers);

// Before next request, wait if needed
await limiter.waitIfNeeded();
```

### Request Interceptors

Add middleware to requests:

```typescript
import { InterceptorManager, BuiltInInterceptors } from '@eclipse-softworks/luna-sdk';

const interceptors = new InterceptorManager();

// Add logging
const { request, response, error } = BuiltInInterceptors.logging(console.log);
interceptors.useRequest(request);
interceptors.useResponse(response);
interceptors.useError(error);

// Add custom headers
interceptors.useRequest(BuiltInInterceptors.customHeaders({
    'X-Custom-Header': 'value',
}));

// Add idempotency keys
interceptors.useRequest(BuiltInInterceptors.idempotencyKey());
```

---

## Webhook Verification

New secure webhook handling:

```typescript
import { WebhookVerifier } from '@eclipse-softworks/luna-sdk';

const webhooks = new WebhookVerifier({
    secret: process.env.LUNA_WEBHOOK_SECRET,
    maxAge: 300, // 5 minutes
});

// In your webhook handler
app.post('/webhooks', express.raw({ type: 'application/json' }), async (req, res) => {
    try {
        const event = await webhooks.verify({
            body: req.body.toString(),
            headers: req.headers,
        });
        
        console.log('Event:', event.payload);
        console.log('Timestamp:', event.timestamp);
        
        res.status(200).send('OK');
    } catch (error) {
        if (error instanceof WebhookVerificationError) {
            console.error('Verification failed:', error.code);
            res.status(400).send('Invalid signature');
        }
    }
});
```

---

## React Hooks (New!)

For React applications:

```bash
npm install @eclipse-softworks/luna-sdk react
```

```tsx
import { LunaProvider, useLuna, useLunaQuery, useLunaMutation } from '@eclipse-softworks/luna-sdk/react';

// Wrap your app
function App() {
    return (
        <LunaProvider config={{ apiKey: process.env.NEXT_PUBLIC_LUNA_API_KEY }}>
            <UserList />
        </LunaProvider>
    );
}

// Use hooks in components
function UserList() {
    const { data, isLoading, error } = useLunaQuery(() => 
        useLuna().users.list({ limit: 10 })
    );

    if (isLoading) return <p>Loading...</p>;
    if (error) return <p>Error: {error.message}</p>;
    
    return (
        <ul>
            {data.items.map(user => (
                <li key={user.id}>{user.name}</li>
            ))}
        </ul>
    );
}
```

---

## Upgrade Steps

1. **Update dependencies:**

```bash
# TypeScript
npm install @eclipse-softworks/luna-sdk@2

# Python
pip install luna-sdk==2.0.0

# Go
go get github.com/eclipse-softworks/luna-sdk-go/v2@latest
```

2. **Update Go imports** (if applicable):

```diff
- import "github.com/eclipse-softworks/luna-sdk-go/luna"
+ import "github.com/eclipse-softworks/luna-sdk-go/v2/luna"
```

3. **Test your application:**

```bash
npm test
```

4. **Update Node.js/Python** if below minimum versions.

---

## Need Help?

- [Full Documentation](https://docs-lunasdk.eclipse-softworks.com)
- [API Reference](https://api.eclipse.dev/docs)
- [GitHub Issues](https://github.com/eclipse-softworks/luna-sdk/issues)
- Email: support@eclipse-softworks.com
