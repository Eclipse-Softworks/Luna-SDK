# Progressive Web Apps (PWA)

Turn your web application into a capable offline-first PWA with Luna SDK.

## Service Worker Setup

Create a service worker file (e.g., `sw.ts`) and use the SDK to generate the necessary code:

```typescript
import { generateServiceWorkerCode } from '@eclipse-softworks/luna-sdk/pwa';

// Generate standard caching and sync logic
const swCode = generateServiceWorkerCode({
    cacheName: 'luna-v1',
    routes: ['/api/*'],
    strategies: {
        '/api/users': 'network-first',
        '/api/static': 'cache-first',
    }
});

// Write this content to your sw.js build output
```

## Client App Setup

Initialize the `PWAManager` in your application entry point:

```typescript
import { PWAManager } from '@eclipse-softworks/luna-sdk/pwa';

const pwa = new PWAManager({
    cacheApi: true,
    backgroundSync: true,
});

// Register the service worker
await pwa.initialize();
```

## Background Sync

Queue requests that fail due to network issues. The service worker will automatically retry them when connectivity returns.

```typescript
try {
    await client.users.create(userData);
} catch (err) {
    if (!navigator.onLine) {
        // Queue for background sync
        await pwa.queueForSync({
            url: '/v1/users',
            method: 'POST',
            body: JSON.stringify(userData)
        });
        alert('Saved offline. Will sync when online.');
    }
}
```

## Install Prompts

Handle the "Add to Home Screen" prompt programmatically:

```typescript
import { InstallPrompt } from '@eclipse-softworks/luna-sdk/pwa';

const prompt = new InstallPrompt();

// Listen for availability
prompt.on('available', () => {
    showInstallButton();
});

// Trigger install
installButton.onClick = async () => {
    const accepted = await prompt.prompt();
    if (accepted) {
        console.log('User installed the app');
    }
};
```
