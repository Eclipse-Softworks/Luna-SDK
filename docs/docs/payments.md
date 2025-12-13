---
sidebar_position: 3
---

# Payments

The Payments module provides unified access to major South African payment gateways.

## Supported Gateways

| Gateway | Features | Best For |
|---------|----------|----------|
| **PayFast** | Cards, EFT, Scan to Pay | E-commerce, Subscriptions |
| **Ozow** | Instant EFT | Real-time bank transfers |
| **Yoco** | Cards, Online | SME / Startup payments |
| **PayShap** | Real-time Clearing | Instant B2B / B2C transfers |

## Usage Examples

### TypeScript

```typescript
const payment = await client.payments.payfast.createPayment({
  amount: 199.99,
  itemName: 'Pro Subscription',
  merchantId: '10000100',
  merchantKey: '46f0cd694581a',
  returnUrl: 'https://example.com/success',
  cancelUrl: 'https://example.com/cancel'
});
console.log('Redirect user to:', payment.url);
```

### Python

```python
payment = await client.payments.payfast.create_payment(
    amount=199.99,
    item_name="Pro Subscription",
    return_url="https://example.com/success"
)
print(f"Redirect user to: {payment.url}")
```

### Go

```go
payment, err := client.Payments().PayFast().CreatePayment(ctx, &luna.PaymentParams{
    Amount:   199.99,
    ItemName: "Pro Subscription",
})
fmt.Printf("Redirect user to: %s\n", payment.URL)
```
