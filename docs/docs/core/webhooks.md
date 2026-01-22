# Webhooks

Listen to events happening in your Luna workspace.

## Supported Events

- `user.created`
- `user.updated`
- `project.created`
- `project.updated`

## Verifying Signatures

```typescript
import { Webhooks } from '@eclipse-softworks/luna-sdk';

const isValid = Webhooks.verifySignature(
  payload,
  signatureHeader,
  webhookSecret
);
```
