---
description: Build and verify all SDK packages
---

# Build All Packages

Build production artifacts for all SDK packages.

// turbo-all

## Steps

1. Build TypeScript SDK:
```bash
cd packages/typescript && npm run build
```

2. Build Python wheel:
```bash
cd packages/python && pip install build && python -m build
```

3. Build Go SDK:
```bash
cd packages/go && go build ./...
```

4. Build CLI for current platform:
```bash
cd cli && go build -o luna ./main.go
```

## Artifacts Created
- `packages/typescript/dist/` - ES modules and type definitions
- `packages/python/dist/` - Wheel and source distribution
- `cli/luna` or `cli/luna.exe` - CLI executable
