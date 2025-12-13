"""Luna SDK - Yoco Integration for Python."""

from __future__ import annotations

import hashlib
import hmac
from datetime import datetime
from typing import Any

from luna.http import HttpClient
from .types import (
    YocoConfig,
    YocoPaymentRequest,
    YocoPayment,
    RefundRequest,
    Refund,
    Amount,
)


class Yoco:
    """Yoco online payment integration."""

    def __init__(self, client: HttpClient, config: YocoConfig) -> None:
        self._client = client
        self._config = config

    async def create_payment(self, request: YocoPaymentRequest) -> YocoPayment:
        """Create a checkout session and get redirect URL."""
        response = await self._client.request(
            method="POST",
            path="/v1/payments/yoco/checkouts",  # Proxy through Luna API
            body={
                "amount": request.amount,
                "currency": request.currency,
                "metadata": request.metadata,
                "successUrl": self._config.success_url if hasattr(self._config, "success_url") else None,
                "cancelUrl": self._config.cancel_url if hasattr(self._config, "cancel_url") else None,
                "failureUrl": self._config.failure_url if hasattr(self._config, "failure_url") else None,
            },
        )
        data = response.data
        return YocoPayment(
            id=data["id"],
            provider="yoco",
            checkout_id=data["checkoutId"],
            amount=Amount(value=data["amount"], currency=data["currency"]),
            status=data["status"],
            reference=data["reference"],
            redirect_url=data["redirectUrl"],
            metadata=data.get("metadata"),
            created_at=data["createdAt"],
            updated_at=data["updatedAt"],
        )

    def verify_webhook(self, payload: str, signature: str) -> bool:
        """Verify webhook signature."""
        expected_signature = hmac.new(
            self._config.secret_key.encode(),
            payload.encode(),
            hashlib.sha256,
        ).hexdigest()
        return signature == expected_signature

    def process_webhook(self, payload: dict[str, Any]) -> YocoPayment:
        """Process webhook event."""
        status_map = {
            "payment.succeeded": "completed",
            "payment.failed": "failed",
            "payment.cancelled": "cancelled",
        }

        event_type = payload.get("type", "")
        payment_data = payload.get("payload", {})

        return YocoPayment(
            id=f"yc_{payment_data.get('id', '')}",
            provider="yoco",
            checkout_id=payment_data.get("id", ""),
            amount=Amount(
                value=payment_data.get("amount", 0),
                currency=payment_data.get("currency", "ZAR"),
            ),
            status=status_map.get(event_type, "pending"),
            reference=payment_data.get("id"),
            metadata=payment_data.get("metadata"),
            redirect_url="",
            created_at=datetime.now().isoformat(),
            updated_at=datetime.now().isoformat(),
        )

    async def refund(self, request: RefundRequest) -> Refund:
        """Request a refund."""
        refund_id = f"ref_{int(datetime.now().timestamp() * 1000)}"

        return Refund(
            id=refund_id,
            payment_id=request.payment_id,
            amount=Amount(value=request.amount or 0, currency="ZAR"),
            status="pending",
            reason=request.reason,
            created_at=datetime.now().isoformat(),
        )
