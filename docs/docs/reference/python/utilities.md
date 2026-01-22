# Python Utilities

The Luna Python SDK includes powerful utilities for South African banking, POPIA compliance, and performance monitoring.

## South African Banking

Validate account numbers and generate compliant EFT references.

```python
from luna.utils import sa_banking

# Validate account number
result = sa_banking.validate_account(
    account_number="1234567890",
    bank_code=sa_banking.BankCode.ABSA
)

if result.is_valid:
    print(f"Valid account for {result.bank.name}")
else:
    print(f"Errors: {result.errors}")

# Generate EFT reference
ref = sa_banking.generate_eft_reference("My Payment Reference 123")
# Output: MYPAYMENTREFERENCE12
```

## POPIA Compliance

Manage user consent and anonymize sensitive data.

```python
from luna.utils import popia

# Create consent record
consent = popia.create_consent_record(
    data_subject="user_123",
    purposes=[popia.ProcessingPurpose.MARKETING],
    expires_in_days=365
)

# Anonymize data
data = {"email": "john@example.com", "id_number": "9001015000087"}
safe_data = popia.anonymize(data, fields=["email"], method="mask")
# Output: {"email": "j***m", "id_number": "9001015000087"}

# Anonymize specific fields
safe_id = popia.anonymize_sa_id(data["id_number"])
# Output: 900101*******
```

## Request Caching

Enable in-memory caching to improve performance.

```python
from luna import LunaClient
from luna.http.cache import CacheConfig

client = LunaClient(
    api_key="lk_prod_xxx",
    cache_config=CacheConfig(enabled=True, ttl=60)
)

# First call hits API
users = client.users.list()

# Second call returns cached data immediately
cached_users = client.users.list()
```

## Metrics Collection

Collect detailed metrics about SDK usage and performance.

```python
from luna.http.metrics import MetricsConfig

def on_metrics(metrics):
    print(f"P95 Latency: {metrics.p95_latency}ms")

client = LunaClient(
    api_key="lk_prod_xxx",
    metrics_config=MetricsConfig(
        enabled=True,
        on_aggregated=on_metrics
    )
)
```
