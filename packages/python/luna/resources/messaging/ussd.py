"""Luna SDK - USSD Integration for Python."""

from __future__ import annotations

from typing import Callable, Awaitable, Optional

from luna.http import HttpClient
from .types import USSDConfig, USSDSession, USSDResponse

# Handler type
USSDHandler = Callable[[USSDSession], Awaitable[USSDResponse] | USSDResponse]


class USSD:
    """USSD service integration for South African networks."""

    def __init__(self, client: HttpClient, config: USSDConfig) -> None:
        self._client = client
        self._config = config
        self._handlers: dict[str, USSDHandler] = {}

    def on_session(self, handler: USSDHandler) -> None:
        """Register a handler for USSD sessions."""
        self._handlers["default"] = handler

    def on_menu(self, path: str, handler: USSDHandler) -> None:
        """Register a handler for a specific menu path."""
        self._handlers[path] = handler

    async def process_request(self, session: USSDSession) -> USSDResponse:
        """Process incoming USSD request."""
        handler = self._handlers.get(session.text) or self._handlers.get("default")

        if not handler:
            return USSDResponse(
                text="Service temporarily unavailable. Please try again later.",
                end=True,
            )

        try:
            result = handler(session)
            if hasattr(result, "__await__"):
                return await result
            return result
        except Exception as e:
            print(f"USSD handler error: {e}")
            return USSDResponse(
                text="An error occurred. Please try again.",
                end=True,
            )

    @staticmethod
    def menu(title: str, options: list[dict[str, str]]) -> str:
        """Create a menu response."""
        lines = [title, ""]
        for option in options:
            lines.append(f"{option['key']}. {option['label']}")
        return "\n".join(lines)

    def parse_africastalking_request(
        self,
        session_id: str,
        phone_number: str,
        service_code: str,
        text: str,
        network_code: Optional[str] = None,
    ) -> USSDSession:
        """Parse Africa's Talking webhook format."""
        from datetime import datetime

        networks = {
            "655001": "Vodacom",
            "655002": "Telkom",
            "655007": "Cell C",
            "655010": "MTN",
        }

        return USSDSession(
            id=f"ussd_{int(datetime.now().timestamp() * 1000)}",
            session_id=session_id,
            phone_number=phone_number,
            service_code=service_code,
            text=text,
            state="active",
            network=networks.get(network_code or "", network_code),
            created_at=datetime.now().isoformat(),
            updated_at=datetime.now().isoformat(),
        )

    def format_africastalking_response(self, response: USSDResponse) -> str:
        """Format response for Africa's Talking."""
        prefix = "END " if response.end else "CON "
        return prefix + response.text

    def parse_clickatell_request(
        self,
        session_id: str,
        msisdn: str,
        request: str,
        shortcode: str,
    ) -> USSDSession:
        """Parse Clickatell webhook format."""
        from datetime import datetime

        return USSDSession(
            id=f"ussd_{int(datetime.now().timestamp() * 1000)}",
            session_id=session_id,
            phone_number=msisdn,
            service_code=shortcode,
            text=request,
            state="active",
            created_at=datetime.now().isoformat(),
            updated_at=datetime.now().isoformat(),
        )

    @property
    def service_code(self) -> str:
        """Get USSD service code."""
        return self._config.service_code

    @staticmethod
    def create_example_menu() -> USSDHandler:
        """Create example session flow."""
        def handler(session: USSDSession) -> USSDResponse:
            parts = [p for p in session.text.split("*") if p]

            if not parts:
                return USSDResponse(
                    text=USSD.menu("Welcome to Luna SDK", [
                        {"key": "1", "label": "Check Balance"},
                        {"key": "2", "label": "Send Payment"},
                        {"key": "3", "label": "Mini Statement"},
                        {"key": "4", "label": "Exit"},
                    ]),
                    end=False,
                )

            selection = parts[0]

            if selection == "1":
                return USSDResponse(text="Your balance is R1,234.56", end=True)
            elif selection == "2":
                if len(parts) == 1:
                    return USSDResponse(text="Enter phone number to send payment:", end=False)
                if len(parts) == 2:
                    return USSDResponse(text="Enter amount (ZAR):", end=False)
                return USSDResponse(text=f"Payment of R{parts[2]} to {parts[1]} initiated.", end=True)
            elif selection == "3":
                return USSDResponse(
                    text="Mini Statement:\n1. Received R500.00\n2. Sent R100.00\n3. Airtime R50.00",
                    end=True,
                )
            elif selection == "4":
                return USSDResponse(text="Thank you for using Luna SDK. Goodbye!", end=True)
            else:
                return USSDResponse(text="Invalid selection. Please try again.", end=True)

        return handler
