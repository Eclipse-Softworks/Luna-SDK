# Desktop Development (Electron)

Build secure, auto-updating desktop applications with Luna SDK's Electron adapter.

## Installation

```bash
npm install @eclipse-softworks/luna-sdk keytar electron-log
```

## Main Process Setup

Initialize the adapter in your `main.ts` to enable secure storage (system keychain) and auto-updates.

```typescript
import { app } from 'electron';
import { ElectronAdapter } from '@eclipse-softworks/luna-sdk/adapters';

const adapter = new ElectronAdapter({
    secureStorage: true,
    autoUpdate: true,
});

app.whenReady().then(async () => {
    // Expose helpers to renderer via IPC
    adapter.setupMainProcessHandlers();
});
```

## Preload Script

Expose safe APIs to the renderer process in `preload.ts`:

```typescript
import { contextBridge } from 'electron';
import { ElectronAdapter } from '@eclipse-softworks/luna-sdk/adapters';

ElectronAdapter.exposeInMainWorld(contextBridge);
```

## Renderer Process

Access SDK features securely from your UI:

```typescript
// Access exposed API
const luna = window.luna;

// Secure Storage
await luna.storage.setApiKey('lk_prod_xxx');
const key = await luna.storage.getApiKey();

// Network Status
luna.onNetworkChange((isOnline) => {
    console.log('Network status:', isOnline);
});
```

## Auto-Updates

The adapter integrates with `electron-updater` to handle checking for updates and installing them automatically.

```typescript
// Check for updates
const updateAvailable = await adapter.checkForUpdates();
```
