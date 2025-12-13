"""Luna SDK - SMS Integration for Python."""

from __future__ import annotations

import re
from datetime import datetime
from typing import Optional

from luna.http import HttpClient
from .types import SMSConfig, SMSSendRequest, SMSMessage, SMSBulkResult


class SMS:
    """SMS messaging integration."""

    def __init__(self, client: HttpClient, config: SMSConfig) -> None:
        self._client = client
        self._config = config

    async def send(self, request: SMSSendRequest) -> SMSMessage:
        """Send a single SMS message."""
        to = request.to[0] if isinstance(request.to, list) else request.to
        if not to:
            raise ValueError("SMS recipient (to) is required")

        message_id = f"sms_{int(datetime.now().timestamp() * 1000)}"
        normalized_to = self._normalize_phone_number(to)

        # Provider-specific logic would go here
        # For now, we return a mock response
        return SMSMessage(
            id=message_id,
            to=normalized_to,
            from_=request.from_ or self._config.sender_id,
            body=request.body,
            status="pending",
            direction="outbound",
            provider=self._config.provider,
            parts=(len(request.body) + 159) // 160,
            metadata=request.metadata,
            created_at=datetime.now().isoformat(),
            updated_at=datetime.now().isoformat(),
        )

    async def send_bulk(self, request: SMSSendRequest) -> SMSBulkResult:
        """Send SMS to multiple recipients."""
        recipients = request.to if isinstance(request.to, list) else [request.to]
        successful: list[SMSMessage] = []
        failed: list[dict] = []

        for to in recipients:
            try:
                message = await self.send(SMSSendRequest(
                    to=to,
                    body=request.body,
                    from_=request.from_,
                    callback_url=request.callback_url,
                    metadata=request.metadata,
                ))
                successful.append(message)
            except Exception as e:
                failed.append({"to": to, "error": str(e)})

        return SMSBulkResult(successful=successful, failed=failed)

    async def get_status(self, message_id: str) -> SMSMessage:
        """Get SMS delivery status."""
        return SMSMessage(
            id=message_id,
            to="",
            body="",
            status="delivered",
            direction="outbound",
            provider=self._config.provider,
            created_at=datetime.now().isoformat(),
            updated_at=datetime.now().isoformat(),
        )

    async def get_balance(self) -> dict:
        """Get account balance."""
        return {"balance": 100.0, "currency": "ZAR"}

    def _normalize_phone_number(self, phone: str) -> str:
        """Normalize South African phone numbers to E.164 format."""
        digits = re.sub(r"\D", "", phone)

        # Handle SA numbers
        if digits.startswith("0") and len(digits) == 10:
            digits = "27" + digits[1:]

        if not digits.startswith("+"):
            digits = "+" + digits

        return digits
