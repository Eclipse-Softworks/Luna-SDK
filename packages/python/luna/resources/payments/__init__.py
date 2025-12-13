"""Luna SDK - Payments Module for South Africa."""

from __future__ import annotations

from typing import Optional

from luna.http import HttpClient
from .types import PaymentsConfig, PayFastConfig, OzowConfig, YocoConfig, PayShapConfig
from .payfast import PayFast
from .ozow import Ozow
from .yoco import Yoco
from .payshap import PayShap

# Re-export types
from .types import (
    Amount,
    PaymentStatus,
    PaymentProvider,
    PayFastPaymentRequest,
    PayFastPayment,
    OzowPaymentRequest,
    OzowPayment,
    YocoPaymentRequest,
    YocoPayment,
    PayShapPaymentRequest,
    PayShapPayment,
    RefundRequest,
    Refund,
)

__all__ = [
    # Main class
    "Payments",
    # Provider classes
    "PayFast",
    "Ozow",
    "Yoco",
    "PayShap",
    # Config types
    "PaymentsConfig",
    "PayFastConfig",
    "OzowConfig",
    "YocoConfig",
    "PayShapConfig",
    # Request/Response types
    "Amount",
    "PaymentStatus",
    "PaymentProvider",
    "PayFastPaymentRequest",
    "PayFastPayment",
    "OzowPaymentRequest",
    "OzowPayment",
    "YocoPaymentRequest",
    "YocoPayment",
    "PayShapPaymentRequest",
    "PayShapPayment",
    "RefundRequest",
    "Refund",
]


class Payments:
    """
    South African Payments resource.

    Supports PayFast, Ozow, Yoco, and PayShap payment gateways.

    Example:
        from luna import LunaClient
        from luna.resources.payments import PayFastConfig

        client = LunaClient(
            api_key="lk_live_xxx",
            payments=PaymentsConfig(
                payfast=PayFastConfig(
                    merchant_id="xxx",
                    merchant_key="xxx",
                ),
            ),
        )

        payment = await client.payments.payfast.create_payment(
            PayFastPaymentRequest(
                amount=199.99,
                item_name="Premium Subscription",
                return_url="https://example.com/success",
                cancel_url="https://example.com/cancel",
                notify_url="https://example.com/webhook",
            )
        )
    """

    def __init__(self, client: HttpClient, config: Optional[PaymentsConfig] = None) -> None:
        self._client = client
        self._config = config or PaymentsConfig()

        # Lazy-initialized providers
        self._payfast: Optional[PayFast] = None
        self._ozow: Optional[Ozow] = None
        self._yoco: Optional[Yoco] = None
        self._payshap: Optional[PayShap] = None

    @property
    def payfast(self) -> PayFast:
        """Get PayFast gateway instance."""
        if not self._payfast:
            if not self._config.payfast:
                raise ValueError(
                    "PayFast not configured. Provide PayFastConfig when initializing LunaClient."
                )
            self._payfast = PayFast(self._client, self._config.payfast)
        return self._payfast

    @property
    def ozow(self) -> Ozow:
        """Get Ozow gateway instance."""
        if not self._ozow:
            if not self._config.ozow:
                raise ValueError(
                    "Ozow not configured. Provide OzowConfig when initializing LunaClient."
                )
            self._ozow = Ozow(self._client, self._config.ozow)
        return self._ozow

    @property
    def yoco(self) -> Yoco:
        """Get Yoco gateway instance."""
        if not self._yoco:
            if not self._config.yoco:
                raise ValueError(
                    "Yoco not configured. Provide YocoConfig when initializing LunaClient."
                )
            self._yoco = Yoco(self._client, self._config.yoco)
        return self._yoco

    @property
    def payshap(self) -> PayShap:
        """Get PayShap gateway instance."""
        if not self._payshap:
            if not self._config.payshap:
                raise ValueError(
                    "PayShap not configured. Provide PayShapConfig when initializing LunaClient."
                )
            self._payshap = PayShap(self._client, self._config.payshap)
        return self._payshap

    def list(self) -> list[str]:
        """List available payment gateways."""
        available = []
        if self._config.payfast:
            available.append("payfast")
        if self._config.ozow:
            available.append("ozow")
        if self._config.yoco:
            available.append("yoco")
        if self._config.payshap:
            available.append("payshap")
        return available
