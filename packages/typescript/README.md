# @eclipse/luna-sdk

Official TypeScript SDK for the Eclipse Softworks Platform API, featuring specialized tools for the South African market.

## Installation

```bash
npm install @eclipse/luna-sdk
```

## Quick Start

```typescript
import { LunaClient } from '@eclipse/luna-sdk';

// Initialize with API Key
// Auto-detects LUNA_API_KEY from environment if not provided
const client = new LunaClient({
  apiKey: process.env.LUNA_API_KEY,
});

async function main() {
  // 1. South African Payments (PayFast, Ozow, Yoco, PayShap)
  const payment = await client.payments.payfast.createPayment({
    amount: 199.99,
    itemName: 'Premium Subscription',
    merchantId: '10000100',
    merchantKey: '46f0cd694581a',
  });
  console.log('Payment URL:', payment.url);

  // 2. Messaging (SMS, WhatsApp)
  await client.messaging.sms.send({
    to: '+27820000000',
    body: 'Your OTP is 12345',
  });

  // 3. South African Business Tools
  // Validate an ID Number
  const idInfo = client.zaTools.idValidation.validate('9001015009087');
  if (idInfo.isValid) {
    console.log(`Valid ID: ${idInfo.gender}, DOB: ${idInfo.dateOfBirth}`);
  }

  // CIPC Company Lookup
  const company = await client.zaTools.cipc.lookup('2020/123456/07');
  console.log('Company:', company.name, company.status);
}

main();
```

## Features

### 🇿🇦 South African Market Ready
*   **Payments**: Native integrations for PayFast, Ozow, Yoco, and PayShap.
*   **Messaging**: SMS (Clickatell/Africa's Talking), WhatsApp Business, USSD.
*   **ZA Tools**: 
    *   CIPC Company Lookup & Verification
    *   B-BBEE Certificate Validation
    *   SA ID Number Validation (Luhn checksum, DOB/Gender extraction)
    *   Address Standardization

### 🔒 Enterprise Grade
*   **POPIA Compliance**: Automatic redaction of sensitive SA data (ID numbers, Tax refs, Banking details) in logs.
*   **Type Safety**: Full TypeScript definitions for all resources.

## Configuration

```typescript
const client = new LunaClient({
  apiKey: 'lk_live_...', 
  // Optional specific configs
  payments: {
    payfast: { merchantId: '...', merchantKey: '...' }
  },
  timeout: 30000, 
  logLevel: 'info', // 'debug' | 'info' | 'warn' | 'error'
});
```

## License

MIT © Eclipse Softworks
