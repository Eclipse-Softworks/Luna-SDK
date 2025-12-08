"""Token-based authentication provider."""

from __future__ import annotations

import asyncio
from datetime import datetime, timedelta
from typing import Awaitable, Callable

import httpx

from luna.auth.types import AuthProvider, TokenPair


class TokenAuth(AuthProvider):
    """
    Token-based authentication provider with automatic refresh.

    Example:
        auth = TokenAuth(
            access_token=session.access_token,
            refresh_token=session.refresh_token,
            on_refresh=save_tokens,
        )
    """

    def __init__(
        self,
        access_token: str,
        refresh_token: str | None = None,
        expires_at: datetime | None = None,
        on_refresh: Callable[[TokenPair], Awaitable[None]] | None = None,
    ) -> None:
        """
        Initialize token authentication.

        Args:
            access_token: OAuth access token
            refresh_token: OAuth refresh token for automatic refresh
            expires_at: Token expiration time
            on_refresh: Callback when tokens are refreshed
        """
        if not access_token:
            raise ValueError("Access token is required")

        self._access_token = access_token
        self._refresh_token = refresh_token
        self._expires_at = expires_at
        self._on_refresh = on_refresh
        self._refresh_lock = asyncio.Lock()

    async def get_headers(self) -> dict[str, str]:
        """Get authorization headers with valid token."""
        if self.needs_refresh():
            await self.refresh()
        return {"Authorization": f"Bearer {self._access_token}"}

    def needs_refresh(self) -> bool:
        """Check if token is expiring soon (within 5 minutes)."""
        if self._expires_at is None:
            return False
        buffer = timedelta(minutes=5)
        return datetime.now() + buffer >= self._expires_at

    async def refresh(self) -> None:
        """Refresh the access token."""
        async with self._refresh_lock:
            # Double-check after acquiring lock
            if not self.needs_refresh():
                return

            if not self._refresh_token:
                raise ValueError("No refresh token available")

            async with httpx.AsyncClient() as client:
                response = await client.post(
                    "https://api.eclipse.dev/v1/auth/refresh",
                    json={"refresh_token": self._refresh_token},
                )
                response.raise_for_status()
                data = response.json()

            self._access_token = data["access_token"]
            self._refresh_token = data["refresh_token"]
            self._expires_at = datetime.now() + timedelta(seconds=data["expires_in"])

            if self._on_refresh:
                token_pair = TokenPair(
                    access_token=self._access_token,
                    refresh_token=self._refresh_token,
                    expires_at=self._expires_at,
                )
                await self._on_refresh(token_pair)

    def update_tokens(self, tokens: TokenPair) -> None:
        """Manually update tokens."""
        self._access_token = tokens.access_token
        self._refresh_token = tokens.refresh_token
        self._expires_at = tokens.expires_at
