---
sidebar_position: 5
---

# South African Tools

Essential business verification and compliance tools for South Africa.

## Available Tools

*   **CIPC**: Company registration lookup and status verification.
*   **B-BBEE**: Compliance level verification.
*   **ID Validation**: South African ID number validation (Luhn checksum, DOB extraction).
*   **Address**: Standardization for SA postal codes and provinces.

## ID Number Validation

The SDK can validate ID numbers **offline** (locally) without making API calls.

### TypeScript
```typescript
const info = client.zaTools.idValidation.validate('9001015009087');
if (info.isValid) {
  console.log(`Citizen: ${info.isSaCitizen}, Gender: ${info.gender}`);
}
```

### Python
```python
info = client.za_tools.validate_id("9001015009087")
if info.is_valid:
    print(f"Born: {info.date_of_birth}")
```

### Go
```go
info := client.ZATools().ValidateID("9001015009087")
```

## CIPC Company Lookup

```python
company = await client.za_tools.cipc.lookup("2020/123456/07")
if company:
    print(f"{company.name} is {company.status}")
```
