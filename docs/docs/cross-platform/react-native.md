# Mobile Development (React Native)

Luna SDK provides first-class support for React Native, including native storage integration, network monitoring, and offline synchronization.

## Installation

```bash
npm install @eclipse-softworks/luna-sdk @react-native-async-storage/async-storage @react-native-community/netinfo
```

## Setup

Wrap your application in the `LunaRNProvider`:

```tsx
import { LunaRNProvider } from '@eclipse-softworks/luna-sdk/adapters';
import { AppRegistry } from 'react-native';
import App from './App';

const Root = () => (
  <LunaRNProvider config={{ apiKey: 'lk_prod_xxx', offlineMode: true }}>
    <App />
  </LunaRNProvider>
);

AppRegistry.registerComponent('MyApp', () => Root);
```

## Features

### Network-Aware Queries

Data fetching that automatically pauses when offline and refetches when connectivity is restored.

```tsx
import { useNetworkAwareQuery } from '@eclipse-softworks/luna-sdk/adapters';
import { Text, View } from 'react-native';

function UserProfile({ userId }) {
  const { data, isOffline } = useNetworkAwareQuery(
    () => client.users.get(userId),
    { refetchOnReconnect: true }
  );

  return (
    <View>
      {isOffline && <Text>You are offline</Text>}
      <Text>User: {data?.name}</Text>
    </View>
  );
}
```

### Offline Mutations

Queue actions (like creating users) when offline, to be automatically retried when back online.

```tsx
import { useOfflineMutation } from '@eclipse-softworks/luna-sdk/adapters';
import { Button } from 'react-native';

function CreateUser() {
  const { mutate, isQueued } = useOfflineMutation(
    (data) => client.users.create(data),
    { queueOffline: true }
  );

  return (
    <Button 
      title="Create User" 
      onPress={() => mutate({ name: 'New User' })} 
    />
  );
}
```

### App State Monitoring

The SDK automatically handles app state changes (background/foreground) to pause/resume synchronization via `useAppState`.
