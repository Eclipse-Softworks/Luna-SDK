---
description: Run all SDK tests (TypeScript, Python, Go)
---

# Run All SDK Tests

Run the full test suite across all SDK packages.

// turbo-all

## Steps

1. Run TypeScript tests:
```bash
cd packages/typescript && npm test
```

2. Run Python tests:
```bash
cd packages/python && python -m pytest
```

3. Run Go tests:
```bash
cd packages/go && go test ./...
```

4. Run CLI build verification:
```bash
cd cli && go build -o luna ./main.go
```

## Expected Results
- All test suites should pass
- CLI should build without errors
