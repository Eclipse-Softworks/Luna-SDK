# Users

Users are the central entities in the Luna platform.

## The User Object

```json
{
  "id": "usr_123",
  "email": "jane@example.com",
  "name": "Jane Doe",
  "avatar_url": "https://...",
  "metadata": {
    "role": "admin"
  },
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

## Operations

### Create a User

```typescript
const user = await client.users.create({
  email: 'jane@example.com',
  name: 'Jane Doe'
});
```

### Get a User

```typescript
const user = await client.users.get('usr_123');
```

### List Users

```typescript
const users = await client.users.list({ limit: 20 });
```

### Update User

```typescript
const user = await client.users.update('usr_123', {
  name: 'Jane Smith'
});
```
