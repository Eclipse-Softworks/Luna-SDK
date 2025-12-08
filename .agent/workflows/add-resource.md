---
description: Add a new resource to all SDKs
---

# Add New Resource

Steps to add a new API resource across all SDK packages.

## 1. Define Types

### TypeScript (`packages/typescript/src/types/index.ts`)
```typescript
export interface NewResource {
  id: string;
  name: string;
  // ... other fields
  created_at: string;
  updated_at: string;
}

export interface NewResourceCreate {
  name: string;
  // ... required fields
}
```

### Python (`packages/python/luna/types/__init__.py`)
```python
class NewResource(BaseModel):
    id: str
    name: str
    created_at: str
    updated_at: str

class NewResourceCreate(BaseModel):
    name: str
```

### Go (`packages/go/luna/resources/types.go`)
```go
type NewResource struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}
```

## 2. Create Resource Class/Struct

### TypeScript (`packages/typescript/src/resources/newresource.ts`)
- Create resource class with CRUD methods
- Export from `packages/typescript/src/resources/index.ts`
- Add to client in `packages/typescript/src/client.ts`

### Python (`packages/python/luna/resources/newresource.py`)
- Create resource class with async CRUD methods
- Export from `packages/python/luna/resources/__init__.py`
- Add to client in `packages/python/luna/client.py`

### Go (`packages/go/luna/resources/newresource.go`)
- Create resource struct with CRUD methods
- Add to client in `packages/go/luna/client.go`

## 3. Add Tests

Create test files for the new resource in each SDK's test directory.

## 4. Update Documentation

Add usage examples to each SDK's README.
