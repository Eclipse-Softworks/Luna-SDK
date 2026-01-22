# Vue 3 & Nuxt Integration

Seamlessly integrate Luna SDK into Vue 3 applications and Nuxt projects.

## Vue 3 (Composables)

Use the `createUseLuna` and `createUseQuery` composables for reactive data fetching.

```vue
<script setup>
import { createUseLuna, createUseQuery } from '@eclipse-softworks/luna-sdk/adapters';

const { client } = createUseLuna({ 
  apiKey: import.meta.env.VITE_LUNA_API_KEY 
});

const { data: users, isLoading, error } = createUseQuery(
  () => client.users.list()
);

const { mutate } = createUseMutation(
  (userData) => client.users.create(userData)
);
</script>

<template>
  <div v-if="isLoading">Loading...</div>
  <div v-else>
    <div v-for="user in users.items" :key="user.id">
      {{ user.name }}
    </div>
    <button @click="mutate({ name: 'New User' })">Add User</button>
  </div>
</template>
```

## Nuxt Module

Register the module in `nuxt.config.ts`:

```typescript
export default defineNuxtConfig({
  modules: ['@eclipse-softworks/luna-sdk/nuxt'],
  luna: {
    apiKey: process.env.LUNA_API_KEY,
    exposeClient: true // Auto-inject $luna
  }
});
```

Usage in components:

```vue
<script setup>
const { $luna } = useNuxtApp();
const { data } = await useAsyncData('users', () => $luna.users.list());
</script>
```
