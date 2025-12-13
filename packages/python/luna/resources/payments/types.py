"""Luna SDK - South African Payments Types."""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import Literal, Optional
from datetime import datetime


# Payment Provider Types
PaymentProvider = Literal["payfast", "ozow", "yoco", "payshap"]
PaymentStatus = Literal["pending", "processing", "completed", "failed", "cancelled", "refunded"]
PaymentMethod = Literal["eft", "cc", "dc", "mp", "cd", "sc"]  # EFT, Credit, Debit, Mobicred, Cash Deposit, SnapScan


@dataclass
class Amount:
    """Monetary amount."""
    value: int  # In cents
    currency: str = "ZAR"


@dataclass
class PayFastConfig:
    """PayFast configuration."""
    merchant_id: str
    merchant_key: str
    passphrase: Optional[str] = None
    sandbox: bool = False


@dataclass
class OzowConfig:
    """Ozow configuration."""
    site_code: str
    private_key: str
    api_key: Optional[str] = None
    sandbox: bool = False


@dataclass
class YocoConfig:
    """Yoco configuration."""
    secret_key: str
    public_key: Optional[str] = None
    sandbox: bool = False


@dataclass
class PayShapConfig:
    """PayShap configuration."""
    merchant_id: str
    bank_id: str
    api_key: Optional[str] = None
    sandbox: bool = False


@dataclass
class PaymentsConfig:
    """Combined payments configuration."""
    payfast: Optional[PayFastConfig] = None
    ozow: Optional[OzowConfig] = None
    yoco: Optional[YocoConfig] = None
    payshap: Optional[PayShapConfig] = None


# Request/Response Models

@dataclass
class PayFastPaymentRequest:
    """PayFast payment request."""
    amount: float
    item_name: str
    return_url: str
    cancel_url: str
    notify_url: str
    item_description: Optional[str] = None
    currency: str = "ZAR"
    email_address: Optional[str] = None
    cell_number: Optional[str] = None
    payment_method: Optional[PaymentMethod] = None
    custom_str1: Optional[str] = None
    custom_str2: Optional[str] = None
    custom_str3: Optional[str] = None
    custom_int1: Optional[int] = None
    custom_int2: Optional[int] = None


@dataclass
class PayFastPayment:
    """PayFast payment response."""
    id: str
    provider: Literal["payfast"] = "payfast"
    amount: Amount = field(default_factory=lambda: Amount(0))
    status: PaymentStatus = "pending"
    reference: Optional[str] = None
    description: Optional[str] = None
    payment_url: str = ""
    signature: Optional[str] = None
    pf_payment_id: Optional[str] = None
    created_at: str = ""
    updated_at: str = ""


@dataclass
class OzowPaymentRequest:
    """Ozow payment request."""
    amount: float
    transaction_reference: str
    bank_reference: str
    success_url: str
    cancel_url: str
    error_url: str
    notify_url: str
    is_test: Optional[bool] = None
    customer_first_name: Optional[str] = None
    customer_last_name: Optional[str] = None
    customer_email: Optional[str] = None
    customer_phone: Optional[str] = None


@dataclass
class OzowPayment:
    """Ozow payment response."""
    id: str
    provider: Literal["ozow"] = "ozow"
    amount: Amount = field(default_factory=lambda: Amount(0))
    status: PaymentStatus = "pending"
    reference: Optional[str] = None
    description: Optional[str] = None
    payment_url: str = ""
    transaction_id: Optional[str] = None
    created_at: str = ""
    updated_at: str = ""


@dataclass
class YocoPaymentRequest:
    """Yoco payment request."""
    amount: int  # In cents
    success_url: str
    cancel_url: str
    failure_url: Optional[str] = None
    currency: str = "ZAR"
    metadata: Optional[dict] = None
    line_items: Optional[list] = None


@dataclass
class YocoPayment:
    """Yoco payment response."""
    id: str
    provider: Literal["yoco"] = "yoco"
    checkout_id: str = ""
    amount: Amount = field(default_factory=lambda: Amount(0))
    status: PaymentStatus = "pending"
    reference: Optional[str] = None
    redirect_url: str = ""
    metadata: Optional[dict] = None
    created_at: str = ""
    updated_at: str = ""


@dataclass
class PayShapPaymentRequest:
    """PayShap payment request."""
    amount: float
    reference: str
    shap_id: Optional[str] = None
    recipient_account: Optional[str] = None
    recipient_bank: Optional[str] = None
    expiry_minutes: int = 30


@dataclass
class PayShapPayment:
    """PayShap payment response."""
    id: str
    provider: Literal["payshap"] = "payshap"
    shap_id: str = ""
    amount: Amount = field(default_factory=lambda: Amount(0))
    status: PaymentStatus = "pending"
    reference: Optional[str] = None
    qr_code: Optional[str] = None
    expires_at: Optional[str] = None
    created_at: str = ""
    updated_at: str = ""


@dataclass
class RefundRequest:
    """Refund request."""
    payment_id: str
    amount: Optional[int] = None  # Partial refund in cents
    reason: Optional[str] = None


@dataclass
class Refund:
    """Refund response."""
    id: str
    payment_id: str
    amount: Amount = field(default_factory=lambda: Amount(0))
    status: Literal["pending", "completed", "failed"] = "pending"
    reason: Optional[str] = None
    created_at: str = ""
