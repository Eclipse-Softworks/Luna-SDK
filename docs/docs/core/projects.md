# Projects

Projects allow you to group resources and manage access.

## The Project Object

```json
{
  "id": "prj_456",
  "name": "My Awesome App",
  "description": "The next big thing",
  "owner_id": "usr_123",
  "metadata": {},
  "created_at": "2024-01-01T00:00:00Z"
}
```

## Operations

### Create Project

```typescript
const project = await client.projects.create({
  name: 'Marketing Campaign',
  owner_id: 'usr_123'
});
```
