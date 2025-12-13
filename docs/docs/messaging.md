---
sidebar_position: 4
---

# Messaging

Unified messaging API for SMS, WhatsApp, and USSD.

## Channels

| Channel | Providers | Features |
|---------|-----------|----------|
| **SMS** | Clickatell, Africa's Talking | OTPs, Notifications |
| **WhatsApp** | WhatsApp Business API | Rich Media, Templates |
| **USSD** | Network Providers | Interactive Menu Sessions |

## Sending SMS

### TypeScript
```typescript
await client.messaging.sms.send({
  to: '+27820000000',
  body: 'Your Luna OTP is 1234'
});
```

### Python
```python
await client.messaging.sms.send(
    to="+27820000000", 
    body="Your Luna OTP is 1234"
)
```

### Go
```go
client.Messaging().SMS().Send(ctx, &luna.SMSParams{
    To:   "+27820000000",
    Body: "Your Luna OTP is 1234",
})
```

## WhatsApp Templates

```python
await client.messaging.whatsapp.send_template(
    to="+27820000000",
    template="welcome_message",
    parameters=["John"]
)
```
