"""Luna SDK - Messaging Types for Python."""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import Literal, Optional
from datetime import datetime


# Provider Types
SMSProvider = Literal["clickatell", "africastalking", "twilio"]
WhatsAppProvider = Literal["cloud_api", "clickatell", "wati", "infobip"]
USSDProvider = Literal["clickatell", "africastalking"]

# Status Types
MessageStatus = Literal["pending", "sent", "delivered", "read", "failed"]
USSDState = Literal["new", "active", "completed", "timeout"]


@dataclass
class SMSConfig:
    """SMS configuration."""
    provider: SMSProvider
    api_key: str
    username: Optional[str] = None
    sender_id: Optional[str] = None
    sandbox: bool = False


@dataclass
class WhatsAppConfig:
    """WhatsApp configuration."""
    provider: WhatsAppProvider
    api_key: str
    phone_number_id: Optional[str] = None
    webhook_token: Optional[str] = None
    sandbox: bool = False


@dataclass
class USSDConfig:
    """USSD configuration."""
    provider: USSDProvider
    api_key: str
    service_code: str
    sandbox: bool = False


@dataclass
class MessagingConfig:
    """Combined messaging configuration."""
    sms: Optional[SMSConfig] = None
    whatsapp: Optional[WhatsAppConfig] = None
    ussd: Optional[USSDConfig] = None


@dataclass
class SMSSendRequest:
    """SMS send request."""
    to: str | list[str]
    body: str
    from_: Optional[str] = None
    callback_url: Optional[str] = None
    metadata: Optional[dict] = None


@dataclass
class SMSMessage:
    """SMS message response."""
    id: str
    to: str
    body: str
    status: MessageStatus = "pending"
    direction: Literal["inbound", "outbound"] = "outbound"
    provider: Optional[SMSProvider] = None
    from_: Optional[str] = None
    parts: int = 1
    metadata: Optional[dict] = None
    created_at: str = ""
    updated_at: str = ""


@dataclass
class SMSBulkResult:
    """SMS bulk send result."""
    successful: list[SMSMessage] = field(default_factory=list)
    failed: list[dict] = field(default_factory=list)


@dataclass
class WhatsAppTextRequest:
    """WhatsApp text message request."""
    to: str
    text: str


@dataclass
class WhatsAppTemplateRequest:
    """WhatsApp template message request."""
    to: str
    template_name: str
    template_params: Optional[dict] = None
    language: str = "en"


@dataclass
class WhatsAppMediaRequest:
    """WhatsApp media message request."""
    to: str
    type: Literal["image", "document", "audio", "video"]
    media_url: str
    caption: Optional[str] = None


@dataclass
class WhatsAppMessage:
    """WhatsApp message response."""
    id: str
    to: str
    type: Literal["text", "template", "image", "document", "audio", "video"] = "text"
    status: MessageStatus = "pending"
    direction: Literal["inbound", "outbound"] = "outbound"
    provider: Optional[WhatsAppProvider] = None
    from_: Optional[str] = None
    text: Optional[str] = None
    template_name: Optional[str] = None
    template_params: Optional[dict] = None
    media_url: Optional[str] = None
    created_at: str = ""
    updated_at: str = ""


@dataclass
class USSDSession:
    """USSD session data."""
    id: str
    session_id: str
    phone_number: str
    service_code: str
    text: str
    state: USSDState = "new"
    network: Optional[str] = None
    created_at: str = ""
    updated_at: str = ""


@dataclass
class USSDResponse:
    """USSD response."""
    text: str
    end: bool = False
