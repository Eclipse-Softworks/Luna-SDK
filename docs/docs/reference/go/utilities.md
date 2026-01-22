# Go Utilities

The Luna Go SDK includes robust utilities for South African banking validation, POPIA compliance, and resilience patterns.

## South African Banking

Located in `github.com/eclipse-softworks/luna-sdk/packages/go/luna/utils`.

```go
import "github.com/eclipse-softworks/luna-sdk/packages/go/luna/utils"

// Validate account
result := utils.ValidateAccount("1234567890", utils.BankAbsa, nil)
if result.IsValid {
    fmt.Printf("Valid account for %s\n", result.Bank.Name)
}

// Generate EFT reference
ref := utils.GenerateEFTReference("Invoice #12345", 20)
// Output: INVOICE12345
```

## POPIA Compliance

Tools regarding the Protection of Personal Information Act.

```go
// Create consent
consent := utils.CreateConsentRecord(utils.CreateConsentRecordOptions{
    DataSubject: "usr_123",
    Purposes:    []string{utils.PurposeMarketing},
    ExpiresInDays: 365,
})

// Check validity
if consent.IsValid() {
    // Process data
}

// Anonymize fields
email := "john@example.com"
masked := utils.AnonymizeEmail(email)
// Output: j***@example.com

id := "9001015000087"
safeID := utils.AnonymizeSAID(id)
// Output: 900101*******
```

## Resilience & Caching

The Go Client includes built-in support for caching, retries, and circuit breakers.

### Caching

```go
import "github.com/eclipse-softworks/luna-sdk/packages/go/luna/http"

config := luna.Config{
    APIKey: "lk_prod_xxx",
    Cache: &http.CacheConfig{
        Enabled: true,
        TTL:     time.Minute,
    },
}
```

### Metrics

```go
collector := http.NewMetricsCollector(&http.MetricsConfig{
    Enabled: true,
    OnAggregatedMetrics: func(m *http.AggregatedMetrics) {
        fmt.Printf("Requests/sec: %.2f\n", m.RequestsPerSecond)
    },
})
```
