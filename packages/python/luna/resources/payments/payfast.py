"""Luna SDK - PayFast Integration for Python."""

from __future__ import annotations

import hashlib
import hmac
from datetime import datetime
from typing import Any
from urllib.parse import urlencode

from luna.http import HttpClient
from .types import (
    PayFastConfig,
    PayFastPaymentRequest,
    PayFastPayment,
    RefundRequest,
    Refund,
    Amount,
)

PAYFAST_LIVE_URL = "https://www.payfast.co.za/eng/process"
PAYFAST_SANDBOX_URL = "https://sandbox.payfast.co.za/eng/process"


class PayFast:
    """PayFast payment gateway integration."""

    def __init__(self, client: HttpClient, config: PayFastConfig) -> None:
        self._client = client
        self._config = config

    async def create_payment(self, request: PayFastPaymentRequest) -> PayFastPayment:
        """Create a payment request and get redirect URL."""
        payment_id = f"pf_{int(datetime.now().timestamp() * 1000)}"

        data: dict[str, str] = {
            "merchant_id": self._config.merchant_id,
            "merchant_key": self._config.merchant_key,
            "return_url": request.return_url,
            "cancel_url": request.cancel_url,
            "notify_url": request.notify_url,
            "m_payment_id": payment_id,
            "amount": f"{request.amount:.2f}",
            "item_name": request.item_name,
        }

        if request.item_description:
            data["item_description"] = request.item_description
        if request.email_address:
            data["email_address"] = request.email_address
        if request.cell_number:
            data["cell_number"] = request.cell_number
        if request.custom_str1:
            data["custom_str1"] = request.custom_str1
        if request.custom_str2:
            data["custom_str2"] = request.custom_str2
        if request.custom_str3:
            data["custom_str3"] = request.custom_str3
        if request.custom_int1:
            data["custom_int1"] = str(request.custom_int1)
        if request.custom_int2:
            data["custom_int2"] = str(request.custom_int2)
        if request.payment_method:
            data["payment_method"] = request.payment_method

        signature = self._generate_signature(data)
        data["signature"] = signature

        base_url = PAYFAST_SANDBOX_URL if self._config.sandbox else PAYFAST_LIVE_URL
        payment_url = f"{base_url}?{urlencode(data)}"

        return PayFastPayment(
            id=payment_id,
            provider="payfast",
            amount=Amount(value=int(request.amount * 100), currency=request.currency),
            status="pending",
            reference=payment_id,
            description=request.item_description,
            payment_url=payment_url,
            signature=signature,
            created_at=datetime.now().isoformat(),
            updated_at=datetime.now().isoformat(),
        )

    def verify_webhook(self, payload: dict[str, Any]) -> bool:
        """Verify webhook signature."""
        signature = payload.pop("signature", "")
        expected_signature = self._generate_signature(payload)
        return signature == expected_signature

    def process_webhook(self, payload: dict[str, Any]) -> PayFastPayment:
        """Process webhook and return payment status."""
        status_map = {
            "COMPLETE": "completed",
            "FAILED": "failed",
            "PENDING": "pending",
            "CANCELLED": "cancelled",
        }

        return PayFastPayment(
            id=payload.get("m_payment_id", ""),
            provider="payfast",
            pf_payment_id=payload.get("pf_payment_id"),
            amount=Amount(
                value=int(float(payload.get("amount_gross", 0)) * 100),
                currency="ZAR",
            ),
            status=status_map.get(payload.get("payment_status", ""), "pending"),
            reference=payload.get("m_payment_id"),
            description=payload.get("item_name"),
            payment_url="",
            signature=payload.get("signature"),
            created_at=datetime.now().isoformat(),
            updated_at=datetime.now().isoformat(),
        )

    async def refund(self, request: RefundRequest) -> Refund:
        """Request a refund for a payment."""
        refund_id = f"ref_{int(datetime.now().timestamp() * 1000)}"

        return Refund(
            id=refund_id,
            payment_id=request.payment_id,
            amount=Amount(value=request.amount or 0, currency="ZAR"),
            status="pending",
            reason=request.reason,
            created_at=datetime.now().isoformat(),
        )

    def _generate_signature(self, data: dict[str, str]) -> str:
        """Generate MD5 signature for PayFast."""
        sorted_keys = sorted(data.keys())
        param_string = "&".join(
            f"{key}={data[key].replace(' ', '+')}"
            for key in sorted_keys
            if data[key] and data[key] != ""
        )

        if self._config.passphrase:
            param_string += f"&passphrase={self._config.passphrase}"

        return hashlib.md5(param_string.encode()).hexdigest()
