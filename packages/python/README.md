# Luna SDK for Python

Official Python SDK for the Eclipse Softworks Platform API, optimized for high performance and South African market integration.

## Installation

```bash
pip install eclipse-luna-sdk
```

## Quick Start

```python
import os
from luna import LunaClient

# AUTO-CONFIGURATION ("Spring Boot" style)
# The client automatically reads LUNA_API_KEY and LUNA_BASE_URL from os.environ
# No need to pass arguments if environment variables are set.
client = LunaClient()

async def main():
    async with client:
        # 1. South African Payments
        payment = await client.payments.payfast.create_payment(
            amount=199.99,
            item_name="Premium Plan",
            return_url="https://example.com/ok"
        )
        print(f"Pay URL: {payment.url}")

        # 2. Messaging
        await client.messaging.whatsapp.send_template(
            to="+27820000000",
            template="welcome_message",
            parameters=["John"]
        )

        # 3. ZA Tools (Strict Mode Example)
        # Validate ID format locally before API call
        id_info = client.za_tools.validate_id("9001015009087")
        if id_info.is_valid:
            print(f"User is {id_info.gender}, born {id_info.date_of_birth}")

        # CIPC Lookup
        company = await client.za_tools.cipc.lookup("2020/123456/07")
        if company:
            print(f"Company: {company.name} ({company.status})")

if __name__ == "__main__":
    import asyncio
    asyncio.run(main())
```

## Features

### 🚀 High-Performance Architecture
*   **Auto-Configuration**: Zero-code initialization using `LUNA_API_KEY` / `LUNA_ACCESS_TOKEN`.
*   **Strict Mode** ("Rust-like Safety"): Enable strict client-side validation for critical SA data (ID checksums, Tax refs) to catch errors before network calls.
*   **Connection Pooling** ("C-like Speed"): Optimized `httpx` transport with tuned keep-alive (`max_keepalive=20`) and connection limits (`max_connections=100`) for high-throughput applications.

### 🇿🇦 South African Market Ready
*   **Payments**: PayFast, Ozow, Yoco, PayShap.
*   **Messaging**: SMS (Clickatell/Africa's Talking), WhatsApp, USSD.
*   **ZA Tools**: CIPC, B-BBEE, ID & Address Validation.

## Configuration

### Strict Mode & Performance

```python
client = LunaClient(
    # Explicit config (optional if env vars set)
    api_key="lk_live_...",
    
    # Enable Strict Validation
    strict=True,
    
    # Connection Config (Optimized defaults are used automatically)
    timeout=30.0,
)
```

## License

MIT © Eclipse Softworks
