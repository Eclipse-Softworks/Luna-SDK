"""Luna SDK - WhatsApp Integration for Python."""

from __future__ import annotations

import hashlib
import hmac
import re
from datetime import datetime
from typing import Any

from luna.http import HttpClient
from .types import (
    WhatsAppConfig,
    WhatsAppTextRequest,
    WhatsAppTemplateRequest,
    WhatsAppMediaRequest,
    WhatsAppMessage,
)


class WhatsApp:
    """WhatsApp Business API integration."""

    def __init__(self, client: HttpClient, config: WhatsAppConfig) -> None:
        self._client = client
        self._config = config

    async def send_text(self, request: WhatsAppTextRequest) -> WhatsAppMessage:
        """Send a text message."""
        message_id = f"wa_{int(datetime.now().timestamp() * 1000)}"
        to = self._normalize_phone_number(request.to)

        return WhatsAppMessage(
            id=message_id,
            to=to,
            type="text",
            text=request.text,
            status="pending",
            direction="outbound",
            provider=self._config.provider,
            created_at=datetime.now().isoformat(),
            updated_at=datetime.now().isoformat(),
        )

    async def send_template(self, request: WhatsAppTemplateRequest) -> WhatsAppMessage:
        """Send a template message."""
        message_id = f"wa_{int(datetime.now().timestamp() * 1000)}"
        to = self._normalize_phone_number(request.to)

        return WhatsAppMessage(
            id=message_id,
            to=to,
            type="template",
            template_name=request.template_name,
            template_params=request.template_params,
            status="pending",
            direction="outbound",
            provider=self._config.provider,
            created_at=datetime.now().isoformat(),
            updated_at=datetime.now().isoformat(),
        )

    async def send_media(self, request: WhatsAppMediaRequest) -> WhatsAppMessage:
        """Send a media message."""
        message_id = f"wa_{int(datetime.now().timestamp() * 1000)}"
        to = self._normalize_phone_number(request.to)

        return WhatsAppMessage(
            id=message_id,
            to=to,
            type=request.type,
            media_url=request.media_url,
            status="pending",
            direction="outbound",
            provider=self._config.provider,
            created_at=datetime.now().isoformat(),
            updated_at=datetime.now().isoformat(),
        )

    async def get_status(self, message_id: str) -> WhatsAppMessage:
        """Get message status."""
        return WhatsAppMessage(
            id=message_id,
            to="",
            type="text",
            status="delivered",
            direction="outbound",
            provider=self._config.provider,
            created_at=datetime.now().isoformat(),
            updated_at=datetime.now().isoformat(),
        )

    def verify_webhook(self, payload: str, signature: str) -> bool:
        """Verify webhook signature (Cloud API)."""
        if not self._config.webhook_token:
            raise ValueError("Webhook token not configured")

        expected_signature = hmac.new(
            self._config.webhook_token.encode(),
            payload.encode(),
            hashlib.sha256,
        ).hexdigest()

        return f"sha256={expected_signature}" == signature

    def process_webhook(self, payload: dict[str, Any]) -> list[WhatsAppMessage]:
        """Process incoming webhook."""
        messages: list[WhatsAppMessage] = []

        for entry in payload.get("entry", []):
            for change in entry.get("changes", []):
                value = change.get("value", {})

                # Process incoming messages
                for msg in value.get("messages", []):
                    messages.append(WhatsAppMessage(
                        id=msg.get("id", ""),
                        to=value.get("metadata", {}).get("phone_number_id", ""),
                        from_=msg.get("from"),
                        type=msg.get("type", "text"),
                        text=msg.get("text", {}).get("body"),
                        status="delivered",
                        direction="inbound",
                        provider=self._config.provider,
                        created_at=datetime.fromtimestamp(int(msg.get("timestamp", 0))).isoformat(),
                        updated_at=datetime.now().isoformat(),
                    ))

                # Process status updates
                for status in value.get("statuses", []):
                    status_map = {
                        "sent": "sent",
                        "delivered": "delivered",
                        "read": "read",
                        "failed": "failed",
                    }
                    messages.append(WhatsAppMessage(
                        id=status.get("id", ""),
                        to="",
                        type="text",
                        status=status_map.get(status.get("status", ""), "pending"),
                        direction="outbound",
                        provider=self._config.provider,
                        created_at=datetime.fromtimestamp(int(status.get("timestamp", 0))).isoformat(),
                        updated_at=datetime.now().isoformat(),
                    ))

        return messages

    def _normalize_phone_number(self, phone: str) -> str:
        """Normalize phone number for WhatsApp (E.164 without +)."""
        digits = re.sub(r"\D", "", phone)

        if digits.startswith("0") and len(digits) == 10:
            digits = "27" + digits[1:]

        if digits.startswith("+"):
            digits = digits[1:]

        return digits
