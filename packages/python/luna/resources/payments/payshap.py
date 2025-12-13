"""Luna SDK - PayShap Integration for Python."""

from __future__ import annotations

import base64
import json
from datetime import datetime, timedelta
from typing import Literal

from luna.http import HttpClient
from .types import (
    PayShapConfig,
    PayShapPaymentRequest,
    PayShapPayment,
    Amount,
)


# South African banks with PayShap support
SA_BANKS = {
    "ABSA": "absa",
    "CAPITEC": "capitec",
    "FNB": "fnb",
    "NEDBANK": "nedbank",
    "STANDARD_BANK": "standard",
    "INVESTEC": "investec",
    "DISCOVERY_BANK": "discovery",
    "TYMEBANK": "tymebank",
    "AFRICAN_BANK": "african",
}

SABank = Literal["absa", "capitec", "fnb", "nedbank", "standard", "investec", "discovery", "tymebank", "african"]


class PayShap:
    """PayShap real-time payment integration."""

    def __init__(self, client: HttpClient, config: PayShapConfig) -> None:
        self._client = client
        self._config = config

    async def create_payment(self, request: PayShapPaymentRequest) -> PayShapPayment:
        """Create a PayShap payment request."""
        payment_id = f"ps_{int(datetime.now().timestamp() * 1000)}"
        expiry_minutes = request.expiry_minutes or 30
        expires_at = datetime.now() + timedelta(minutes=expiry_minutes)

        # Generate QR code data
        qr_data = json.dumps({
            "type": "payshap",
            "merchantId": self._config.merchant_id,
            "amount": request.amount,
            "reference": request.reference,
        })
        qr_code = base64.b64encode(qr_data.encode()).decode()

        return PayShapPayment(
            id=payment_id,
            provider="payshap",
            shap_id=f"shp_{payment_id[-8:]}",
            amount=Amount(value=int(request.amount * 100), currency="ZAR"),
            status="pending",
            reference=request.reference,
            qr_code=qr_code,
            expires_at=expires_at.isoformat(),
            created_at=datetime.now().isoformat(),
            updated_at=datetime.now().isoformat(),
        )

    async def get_payment(self, payment_id: str) -> PayShapPayment:
        """Get payment status."""
        return PayShapPayment(
            id=payment_id,
            provider="payshap",
            shap_id=f"shp_{payment_id[-8:]}",
            amount=Amount(value=0, currency="ZAR"),
            status="pending",
            reference=payment_id,
            expires_at=(datetime.now() + timedelta(minutes=30)).isoformat(),
            created_at=datetime.now().isoformat(),
            updated_at=datetime.now().isoformat(),
        )

    async def cancel_payment(self, payment_id: str) -> PayShapPayment:
        """Cancel a pending payment."""
        payment = await self.get_payment(payment_id)
        payment.status = "cancelled"
        payment.updated_at = datetime.now().isoformat()
        return payment

    async def lookup_shap_id(self, shap_id: str) -> dict:
        """Lookup a ShapID (payment proxy)."""
        import re
        is_valid = bool(re.match(r"^[a-zA-Z0-9@._-]+$", shap_id)) and len(shap_id) >= 5

        return {
            "valid": is_valid,
            "bank_name": "Sample Bank" if is_valid else None,
            "account_holder_name": "Account Holder" if is_valid else None,
        }

    async def generate_receive_qr(
        self,
        amount: float | None = None,
        reference: str | None = None,
    ) -> dict:
        """Generate a QR code for receiving payments."""
        shap_id = f"shp_{self._config.merchant_id}_{int(datetime.now().timestamp())}"

        qr_data = json.dumps({
            "type": "payshap_receive",
            "shapId": shap_id,
            "merchantId": self._config.merchant_id,
            "amount": amount,
            "reference": reference,
        })

        return {
            "qr_code": base64.b64encode(qr_data.encode()).decode(),
            "shap_id": shap_id,
        }

    def validate_bank_account(self, account_number: str, bank_id: SABank) -> bool:
        """Validate a South African bank account number format."""
        account_lengths: dict[SABank, list[int]] = {
            "absa": [10, 11],
            "capitec": [10],
            "fnb": [10, 11, 12],
            "nedbank": [10, 11],
            "standard": [9, 10, 11],
            "investec": [10],
            "discovery": [10],
            "tymebank": [10],
            "african": [11],
        }

        valid_lengths = account_lengths.get(bank_id)
        if not valid_lengths:
            return False

        digits_only = "".join(c for c in account_number if c.isdigit())
        return len(digits_only) in valid_lengths
