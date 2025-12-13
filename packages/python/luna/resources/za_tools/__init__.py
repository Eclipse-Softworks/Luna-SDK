"""Luna SDK - South African Business Tools for Python."""

from __future__ import annotations

from typing import Optional

from luna.http import HttpClient
from .types import ZAToolsConfig, CIPCConfig, BBBEEConfig
from .cipc import CIPC
from .bbbee import BBBEE
from .id_validation import IDValidation
from .address import AddressUtils

# Re-export types
from .types import (
    CompanyType,
    CompanyStatus,
    Company,
    Director,
    BBBEELevel,
    BBBEECertificate,
    BBBEEScorecard,
    SAIDInfo,
    SAProvince,
    SAAddress,
)

__all__ = [
    # Main class
    "ZATools",
    # Provider classes
    "CIPC",
    "BBBEE",
    "IDValidation",
    "AddressUtils",
    # Config types
    "ZAToolsConfig",
    "CIPCConfig",
    "BBBEEConfig",
    # Types
    "CompanyType",
    "CompanyStatus",
    "Company",
    "Director",
    "BBBEELevel",
    "BBBEECertificate",
    "BBBEEScorecard",
    "SAIDInfo",
    "SAProvince",
    "SAAddress",
]


class ZATools:
    """
    South African Business Tools.

    Provides utilities for:
    - CIPC company lookup and verification
    - B-BBEE compliance verification
    - SA ID number validation
    - Address validation and formatting

    Example:
        from luna import LunaClient
        from luna.resources.za_tools import ZAToolsConfig

        client = LunaClient(
            api_key="lk_live_xxx",
            za_tools=ZAToolsConfig(cipc=CIPCConfig()),
        )

        # Validate SA ID
        id_info = client.za_tools.id_validation.validate("9001015009087")
        print(id_info.is_valid, id_info.gender, id_info.date_of_birth)

        # Look up company
        company = await client.za_tools.cipc.lookup("2020/123456/07")
    """

    def __init__(self, client: HttpClient, config: Optional[ZAToolsConfig] = None) -> None:
        self._client = client
        self._config = config or ZAToolsConfig()

        # Lazy-initialized services
        self._cipc: Optional[CIPC] = None
        self._bbbee: Optional[BBBEE] = None

        # Standalone utilities (no API calls needed)
        self._id_validation: Optional[IDValidation] = None
        self._address: Optional[AddressUtils] = None

    @property
    def cipc(self) -> CIPC:
        """Get CIPC service instance."""
        if not self._cipc:
            self._cipc = CIPC(
                self._client, 
                self._config.cipc or CIPCConfig(), 
                strict=self._config.strict
            )
        return self._cipc

    @property
    def bbbee(self) -> BBBEE:
        """Get B-BBEE service instance."""
        if not self._bbbee:
            self._bbbee = BBBEE(self._client, self._config.bbbee or BBBEEConfig())
        return self._bbbee

    @property
    def id_validation(self) -> IDValidation:
        """Get ID validation utility (no API needed)."""
        if not self._id_validation:
            self._id_validation = IDValidation()
        return self._id_validation

    @property
    def address(self) -> AddressUtils:
        """Get address utilities (no API needed)."""
        if not self._address:
            self._address = AddressUtils()
        return self._address

    def validate_id(self, id_number: str) -> SAIDInfo:
        """Convenience method to validate SA ID number."""
        return self.id_validation.validate(id_number)

    def list(self) -> list[str]:
        """List available ZA tools."""
        return ["cipc", "bbbee", "id_validation", "address"]
