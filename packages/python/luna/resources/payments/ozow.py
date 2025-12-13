"""Luna SDK - Ozow Integration for Python."""

from __future__ import annotations

import hashlib
from datetime import datetime
from typing import Any
from urllib.parse import urlencode

from luna.http import HttpClient
from .types import (
    OzowConfig,
    OzowPaymentRequest,
    OzowPayment,
    RefundRequest,
    Refund,
    Amount,
)

OZOW_PAYMENT_URL = "https://pay.ozow.com"


class Ozow:
    """Ozow instant EFT payment integration."""

    def __init__(self, client: HttpClient, config: OzowConfig) -> None:
        self._client = client
        self._config = config

    async def create_payment(self, request: OzowPaymentRequest) -> OzowPayment:
        """Create a payment request and get redirect URL."""
        payment_id = f"oz_{int(datetime.now().timestamp() * 1000)}"

        data: dict[str, str] = {
            "SiteCode": self._config.site_code,
            "CountryCode": "ZA",
            "CurrencyCode": "ZAR",
            "Amount": f"{request.amount:.2f}",
            "TransactionReference": request.transaction_reference,
            "BankReference": request.bank_reference,
            "CancelUrl": request.cancel_url,
            "ErrorUrl": request.error_url,
            "SuccessUrl": request.success_url,
            "NotifyUrl": request.notify_url,
            "IsTest": "true" if (request.is_test or self._config.sandbox) else "false",
        }

        if request.customer_first_name:
            data["CustomerFirstName"] = request.customer_first_name
        if request.customer_last_name:
            data["CustomerLastName"] = request.customer_last_name
        if request.customer_email:
            data["CustomerEmail"] = request.customer_email
        if request.customer_phone:
            data["CustomerPhone"] = request.customer_phone

        hash_string = self._generate_hash_string(data)
        hash_check = self._generate_hash(hash_string)
        data["HashCheck"] = hash_check

        payment_url = f"{OZOW_PAYMENT_URL}?{urlencode(data)}"

        return OzowPayment(
            id=payment_id,
            provider="ozow",
            amount=Amount(value=int(request.amount * 100), currency="ZAR"),
            status="pending",
            reference=request.transaction_reference,
            description=request.bank_reference,
            payment_url=payment_url,
            transaction_id=payment_id,
            created_at=datetime.now().isoformat(),
            updated_at=datetime.now().isoformat(),
        )

    def verify_webhook(self, payload: dict[str, Any]) -> bool:
        """Verify webhook hash."""
        received_hash = payload.pop("Hash", "")
        hash_string = self._generate_hash_string(payload)
        expected_hash = self._generate_hash(hash_string)
        return received_hash.lower() == expected_hash.lower()

    def process_webhook(self, payload: dict[str, Any]) -> OzowPayment:
        """Process webhook and return payment status."""
        status_map = {
            "Complete": "completed",
            "Cancelled": "cancelled",
            "Error": "failed",
            "Abandoned": "cancelled",
            "PendingInvestigation": "processing",
        }

        return OzowPayment(
            id=payload.get("TransactionReference", ""),
            provider="ozow",
            transaction_id=payload.get("TransactionId"),
            amount=Amount(
                value=int(float(payload.get("Amount", 0)) * 100),
                currency="ZAR",
            ),
            status=status_map.get(payload.get("Status", ""), "pending"),
            reference=payload.get("TransactionReference"),
            payment_url="",
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

    def _generate_hash_string(self, data: dict[str, str]) -> str:
        """Generate hash string for Ozow."""
        ordered_fields = [
            "SiteCode", "CountryCode", "CurrencyCode", "Amount",
            "TransactionReference", "BankReference", "CancelUrl",
            "ErrorUrl", "SuccessUrl", "NotifyUrl", "IsTest",
        ]
        return "".join(data.get(field, "") for field in ordered_fields).lower()

    def _generate_hash(self, input_string: str) -> str:
        """Generate SHA512 hash."""
        to_hash = input_string + self._config.private_key.lower()
        return hashlib.sha512(to_hash.encode()).hexdigest()
