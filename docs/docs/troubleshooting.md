# Troubleshooting Guide

Common issues and solutions when using Luna SDK.

## Connection Issues

### Network Timeout Errors

**Symptom:** `NetworkError: Request timeout` or `NETWORK_TIMEOUT` error code.

**Solutions:**

1. **Increase timeout:**
```typescript
const client = new LunaClient({
    apiKey: process.env.LUNA_API_KEY,
    timeout: 60000, // 60 seconds
});
```

2. **Check network connectivity:**
```bash
# Test connection to Luna API
curl -I https://api.eclipse.dev/health
```

3. **Check firewall settings:** Ensure outbound HTTPS (port 443) to `api.eclipse.dev` is allowed.

### Connection Refused

**Symptom:** `NetworkError: Connection error` or `NETWORK_CONNECTION` error code.

**Causes:**
- No internet connection
- DNS resolution failure
- Proxy configuration issues

**Solutions:**

1. Verify DNS resolution:
```bash
nslookup api.eclipse.dev
```

2. If behind a proxy, configure it:
```typescript
// Use environment variables
process.env.HTTPS_PROXY = 'http://proxy.company.com:8080';
```

---

## Authentication Issues

### Invalid API Key

**Symptom:** `AuthenticationError` with code `INVALID_API_KEY`.

**Solutions:**

1. **Verify key format:** Keys should match `lk_<env>_<key>` pattern:
   - `lk_prod_xxx...` for production
   - `lk_test_xxx...` for testing
   - `lk_dev_xxx...` for development

2. **Check environment variable:**
```bash
# Verify the key is set
echo $LUNA_API_KEY
```

3. **Regenerate key:** If compromised, regenerate from [Eclipse Developer Portal](https://developer.eclipse.dev).

### Token Expired

**Symptom:** `AuthenticationError` with code `TOKEN_EXPIRED`.

**Solutions:**

1. **Enable automatic refresh:**
```typescript
const client = new LunaClient({
    accessToken: token,
    refreshToken: refreshToken,
    onTokenRefresh: async (tokens) => {
        // Save new tokens
        await saveTokens(tokens);
    },
});
```

2. **Re-authenticate:** If refresh token is also expired, trigger a new login flow.

---

## Rate Limiting

### Rate Limit Exceeded

**Symptom:** `RateLimitError` with `retryAfter` value.

**Solutions:**

1. **Wait and retry:**
```typescript
try {
    await client.users.list();
} catch (error) {
    if (error instanceof RateLimitError) {
        console.log(`Rate limited. Retry after ${error.retryAfter}s`);
        await sleep(error.retryAfter * 1000);
        // Retry
    }
}
```

2. **Use built-in rate limiter:**
```typescript
import { RateLimiter } from '@eclipse-softworks/luna-sdk';

const limiter = new RateLimiter({ maxWaitTime: 30000 });
await limiter.waitIfNeeded();
```

3. **Request higher limits:** Contact support for enterprise rate limits.

---

## Debug Mode

Enable verbose logging to diagnose issues:

```typescript
const client = new LunaClient({
    apiKey: process.env.LUNA_API_KEY,
    logLevel: 'debug',
});
```

### Custom Logger

```typescript
const client = new LunaClient({
    apiKey: process.env.LUNA_API_KEY,
    logger: {
        debug: (msg, data) => console.debug(`[LUNA] ${msg}`, data),
        info: (msg, data) => console.info(`[LUNA] ${msg}`, data),
        warn: (msg, data) => console.warn(`[LUNA] ${msg}`, data),
        error: (msg, data) => console.error(`[LUNA] ${msg}`, data),
    },
});
```

---

## TypeScript Issues

### Type Errors After Update

**Symptom:** TypeScript compilation errors after SDK update.

**Solutions:**

1. **Update TypeScript:**
```bash
npm install typescript@latest
```

2. **Clear caches:**
```bash
rm -rf node_modules/.cache
rm -rf dist
```

3. **Regenerate types:**
```bash
npm run typecheck
```

### Module Resolution Errors

**Symptom:** `Cannot find module '@eclipse-softworks/luna-sdk'`

**Solutions:**

1. **Check Node.js version:** Requires Node.js 18+
```bash
node --version
```

2. **Verify installation:**
```bash
npm ls @eclipse-softworks/luna-sdk
```

3. **ESM vs CJS:** If using CommonJS, ensure proper imports:
```javascript
// CommonJS
const { LunaClient } = require('@eclipse-softworks/luna-sdk');

// ESM
import { LunaClient } from '@eclipse-softworks/luna-sdk';
```

---

## Webhook Issues

### Signature Verification Failed

**Symptom:** `WebhookVerificationError` with code `INVALID_SIGNATURE`.

**Causes:**
- Wrong webhook secret
- Request body was modified (parsed as JSON)
- Clock skew between servers

**Solutions:**

1. **Use raw body:**
```typescript
// Express.js
app.post('/webhooks', express.raw({ type: 'application/json' }), (req, res) => {
    const event = await webhooks.verify({
        body: req.body.toString(), // Raw string, not parsed
        headers: req.headers,
    });
});
```

2. **Verify secret:** Check webhook secret in Luna Dashboard matches your code.

3. **Check server time:** Ensure server clock is synchronized (use NTP).

---

## Getting Help

If issues persist:

1. **Search existing issues:** [GitHub Issues](https://github.com/eclipse-softworks/luna-sdk/issues)
2. **Open a new issue:** Include SDK version, Node.js version, and error details
3. **Contact support:** support@eclipse-softworks.com
