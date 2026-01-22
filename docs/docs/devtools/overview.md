# DevTools Panel

Debug and monitor your Luna SDK integration with the built-in DevTools overlay.

## Enabling DevTools

The DevTools panel is included in the SDK but disabled by default. Enable it during development:

```typescript
import { DevTools } from '@eclipse-softworks/luna-sdk/devtools';

if (process.env.NODE_ENV === 'development') {
    DevTools.enable({ 
        overlay: true,      // Show floating button
        logRequests: true,  // Log all requests to console
        trackErrors: true   // Monitor SDK errors
    });
}
```

Once enabled, you can toggle the panel with **Ctrl+Shift+L** or by clicking the floating Luna button.

## Request Recorder

Capture sequences of API interactions to replay them later (useful for debugging or creating reproduction scripts).

```typescript
import { RequestRecorder } from '@eclipse-softworks/luna-sdk/devtools';

const recorder = new RequestRecorder();

// Start recording
recorder.start();

// Perform actions in your app...
await client.users.list();
await client.projects.create({ name: 'Demo' });

// Stop and save
const session = recorder.stop();
console.log(JSON.stringify(session, null, 2));

// Replay a session
await recorder.replay(session, client);
```

## Network Statistics

The DevTools panel provides real-time metrics:
- **Request Latency**: P50, P95, P99
- **Error Rates**: By endpoint and error code
- **Cache Hit Rate**: Efficiency of client-side caching
