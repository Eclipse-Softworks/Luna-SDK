"""Luna SDK - Messaging Module for Python."""

from __future__ import annotations

from typing import Optional

from luna.http import HttpClient
from .types import MessagingConfig, SMSConfig, WhatsAppConfig, USSDConfig
from .sms import SMS
from .whatsapp import WhatsApp
from .ussd import USSD

# Re-export types
from .types import (
    SMSProvider,
    WhatsAppProvider,
    USSDProvider,
    MessageStatus,
    USSDState,
    SMSSendRequest,
    SMSMessage,
    SMSBulkResult,
    WhatsAppTextRequest,
    WhatsAppTemplateRequest,
    WhatsAppMediaRequest,
    WhatsAppMessage,
    USSDSession,
    USSDResponse,
)

__all__ = [
    # Main class
    "Messaging",
    # Provider classes
    "SMS",
    "WhatsApp",
    "USSD",
    # Config types
    "MessagingConfig",
    "SMSConfig",
    "WhatsAppConfig",
    "USSDConfig",
    # Types
    "SMSProvider",
    "WhatsAppProvider",
    "USSDProvider",
    "MessageStatus",
    "USSDState",
    "SMSSendRequest",
    "SMSMessage",
    "SMSBulkResult",
    "WhatsAppTextRequest",
    "WhatsAppTemplateRequest",
    "WhatsAppMediaRequest",
    "WhatsAppMessage",
    "USSDSession",
    "USSDResponse",
]


class Messaging:
    """
    Messaging resource for SMS, WhatsApp, and USSD.

    Example:
        from luna import LunaClient
        from luna.resources.messaging import SMSConfig

        client = LunaClient(
            api_key="lk_live_xxx",
            messaging=MessagingConfig(
                sms=SMSConfig(
                    provider="clickatell",
                    api_key="xxx",
                ),
            ),
        )

        message = await client.messaging.sms.send(
            SMSSendRequest(
                to="+27821234567",
                body="Your OTP is 123456",
            )
        )
    """

    def __init__(self, client: HttpClient, config: Optional[MessagingConfig] = None) -> None:
        self._client = client
        self._config = config or MessagingConfig()

        # Lazy-initialized providers
        self._sms: Optional[SMS] = None
        self._whatsapp: Optional[WhatsApp] = None
        self._ussd: Optional[USSD] = None

    @property
    def sms(self) -> SMS:
        """Get SMS instance."""
        if not self._sms:
            if not self._config.sms:
                raise ValueError(
                    "SMS not configured. Provide SMSConfig when initializing LunaClient."
                )
            self._sms = SMS(self._client, self._config.sms)
        return self._sms

    @property
    def whatsapp(self) -> WhatsApp:
        """Get WhatsApp instance."""
        if not self._whatsapp:
            if not self._config.whatsapp:
                raise ValueError(
                    "WhatsApp not configured. Provide WhatsAppConfig when initializing LunaClient."
                )
            self._whatsapp = WhatsApp(self._client, self._config.whatsapp)
        return self._whatsapp

    @property
    def ussd(self) -> USSD:
        """Get USSD instance."""
        if not self._ussd:
            if not self._config.ussd:
                raise ValueError(
                    "USSD not configured. Provide USSDConfig when initializing LunaClient."
                )
            self._ussd = USSD(self._client, self._config.ussd)
        return self._ussd

    def list(self) -> list[str]:
        """List available messaging channels."""
        available = []
        if self._config.sms:
            available.append("sms")
        if self._config.whatsapp:
            available.append("whatsapp")
        if self._config.ussd:
            available.append("ussd")
        return available
