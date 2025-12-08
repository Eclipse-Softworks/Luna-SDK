---
description: Lint and format all code
---

# Lint All Code

Run linters and formatters across all packages.

// turbo-all

## Steps

1. Lint TypeScript:
```bash
cd packages/typescript && npm run lint
```

2. Typecheck TypeScript:
```bash
cd packages/typescript && npm run typecheck
```

3. Lint Python with Ruff:
```bash
cd packages/python && ruff check .
```

4. Format Python with Ruff:
```bash
cd packages/python && ruff format --check .
```

5. Lint Go with golangci-lint:
```bash
cd packages/go && golangci-lint run
```

## Auto-fix Options

To auto-fix issues:

- TypeScript: `npm run lint -- --fix`
- Python: `ruff check --fix .` and `ruff format .`
- Go: `golangci-lint run --fix`
