# Luna SDK

> **Eclipse Softworks** — Cross-Language SDK for Platform Services

[![TypeScript](https://img.shields.io/badge/TypeScript-v1.0.1-blue)](./packages/typescript)
[![Python](https://img.shields.io/badge/Python-v1.0.1-blue)](./packages/python)
[![Go](https://img.shields.io/badge/Go-v1.0.1-blue)](./packages/go)

## Overview
Luna SDK provides official client libraries for the Eclipse Softworks Platform API. Available in TypeScript, Python, and Go with consistent patterns across all languages.
<img width="1462" height="6752" alt="py_basic" src="https://github.com/user-attachments/assets/6cd6d8d9-9217-46c7-a510-db3026ea9bed" />


## Installation

### TypeScript
```bash
npm install @eclipse-softworks/luna-sdk
```

### Python
```bash
pip install luna-sdk
```

### Go
```bash
go get github.com/eclipse-softworks/luna-sdk-go
```

## Quick Start

### TypeScript
```typescript
import { LunaClient } from '@eclipse-softworks/luna-sdk';

const client = new LunaClient({
  apiKey: process.env.LUNA_API_KEY,
});

const users = await client.users.list();
```

### Python
```python
from luna import LunaClient

client = LunaClient(api_key=os.environ["LUNA_API_KEY"])

users = await client.users.list()
```

### Go
```go
import "github.com/eclipse-softworks/luna-sdk-go/luna"

client := luna.NewClient(luna.WithAPIKey(os.Getenv("LUNA_API_KEY")))

users, err := client.Users().List(ctx)
```

## Documentation

Full documentation is available at **[docs-lunasdk.eclipse-softworks.com](https://docs-lunasdk.eclipse-softworks.com)**.

- [Getting Started](https://docs-lunasdk.eclipse-softworks.com/docs/intro)
- [Authentication](https://docs-lunasdk.eclipse-softworks.com/docs/authentication)
- [Service Modules](https://docs-lunasdk.eclipse-softworks.com/docs/services)
- [Error Reference](https://docs-lunasdk.eclipse-softworks.com/docs/errors)

## CLI

```bash
# Install
go install github.com/eclipse-softworks/luna-sdk/cli@latest

# Usage
luna auth login
luna users list
```

## License

MIT © Eclipse Softworks
