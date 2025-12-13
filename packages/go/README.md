# Luna SDK for Go

Official Go SDK for the Eclipse Softworks Platform API, engineered for performance and reliability with deep South African market integrations.

## Installation

```bash
go get github.com/eclipse-softworks/luna-sdk-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/eclipse-softworks/luna-sdk-go/luna"
)

func main() {
    // AUTO-CONFIGURATION ("Spring Boot" style)
    // Automatically loads LUNA_API_KEY from environment variables.
    // Use luna.WithStrictMode() for "Rust-like" safety.
    client, err := luna.NewClient(
        luna.WithStrictMode(),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // 1. South African Payments
    payment, err := client.Payments().PayFast().CreatePayment(ctx, &luna.PaymentParams{
        Amount:   199.99,
        ItemName: "Pro Subscription",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Pay URL: %s\n", payment.URL)

    // 2. Messaging
    msg, _ := client.Messaging().SMS().Send(ctx, &luna.SMSParams{
        To:   "+27820000000",
        Body: "Your OTP is 9988",
    })
    fmt.Println("SMS Sent:", msg.ID)

    // 3. ZA Tools (Strict Mode Active)
    // This will error locally if ID is invalid, saving an API call.
    idInfo := client.ZATools().ValidateID("9001015009087")
    if idInfo.IsValid {
        fmt.Printf("Valid ID: %s (Born: %s)\n", idInfo.Gender, idInfo.DateOfBirth)
    }

    // CIPC Lookup
    company, _ := client.ZATools().CIPC().Lookup(ctx, "2020/123456/07")
    fmt.Printf("Company: %s [%s]\n", company.Name, company.Status)
}
```

## Features

### 🚀 High-Performance Architecture
*   **Auto-Configuration ("Spring Boot")**: Zero-code setup. Reads `LUNA_API_KEY`, `LUNA_BASE_URL` from env.
*   **Strict Mode ("Rust")**: Opt-in strict validation (`luna.WithStrictMode()`) enforces regex checks on client side for ID numbers/Reg numbers.
*   **Speed ("C")**: Highly tuned `http.Transport` with connection pooling (`MaxIdleConns: 100`) for high-concurrency workloads.

### 🇿🇦 South African Market Ready
*   **Payments**: PayFast, Ozow, Yoco, PayShap.
*   **Messaging**: SMS, WhatsApp, USSD.
*   **ZA Tools**: CIPC, B-BBEE, ID/Address Validation.

## Authentication

### API Key

```go
client := luna.NewClient(luna.WithAPIKey("lk_live_xxxx"))
```

### OAuth Token

```go
client := luna.NewClient(luna.WithTokens(accessToken, refreshToken))
```

## Error Handling

```go
user, err := client.Users().Get(ctx, "usr_nonexistent")
if err != nil {
    switch e := err.(type) {
    case *luna.NotFoundError:
        fmt.Printf("User not found: %s\n", e.Message)
    case *luna.RateLimitError:
        fmt.Printf("Rate limited, retry after: %d\n", e.RetryAfter)
    default:
        log.Fatal(err)
    }
}
```

## Configuration

```go
client := luna.NewClient(
    luna.WithAPIKey(os.Getenv("LUNA_API_KEY")),
    luna.WithBaseURL("https://api.staging.eclipse.dev"),
    luna.WithTimeout(60000),
    luna.WithMaxRetries(5),
    luna.WithLogLevel(luna.LogLevelDebug),
)
```

## License

MIT © Eclipse Softworks
