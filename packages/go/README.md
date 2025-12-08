# luna-sdk-go

Official Go SDK for the Eclipse Softworks Platform API.

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
    // API Key authentication
    client := luna.NewClient(
        luna.WithAPIKey(os.Getenv("LUNA_API_KEY")),
    )

    ctx := context.Background()

    // List users
    users, err := client.Users().List(ctx, &luna.ListParams{Limit: 10})
    if err != nil {
        log.Fatal(err)
    }

    for _, user := range users.Data {
        fmt.Printf("User: %s (%s)\n", user.Name, user.Email)
    }

    // Get a specific user
    user, err := client.Users().Get(ctx, "usr_123")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Got user: %s\n", user.Name)
}
```

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
