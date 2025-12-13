"""Luna SDK - CIPC Company Lookup for Python."""

from __future__ import annotations

import re
from typing import Optional

from luna.http import HttpClient
from .types import CIPCConfig, Company, CompanyType, Director


class CIPC:
    """Companies and Intellectual Property Commission integration."""

    def __init__(self, client: HttpClient, config: CIPCConfig, strict: bool = False) -> None:
        self._client = client
        self._config = config
        self._strict = strict

    async def lookup(self, registration_number: str) -> Company | None:
        """Search for a company by registration number."""
        cleaned = re.sub(r"[\s/]", "", registration_number)

        # Strict Validation ("Rust" safety)
        if self._strict and not self.is_valid_registration_number(cleaned):
            raise ValueError(f"Invalid registration number format (strict mode): {registration_number}")

        try:
            response = await self._client.request(
                method="GET",
                path=f"/v1/za/cipc/companies/{cleaned}",
            )
            return Company(**response.data)
        except Exception:
            return None

    async def search_by_name(self, name: str, limit: int = 10) -> list[Company]:
        """Search for companies by name."""
        response = await self._client.request(
            method="GET",
            path="/v1/za/cipc/companies",
            query={"name": name, "limit": str(limit)},
        )
        return [Company(**c) for c in response.data.get("data", [])]

    async def verify(self, registration_number: str) -> dict:
        """Verify if a company is registered and active."""
        company = await self.lookup(registration_number)

        if not company:
            return {"exists": False, "is_active": False}

        return {
            "exists": True,
            "is_active": company.status == "active",
            "company": company,
        }

    async def check_name_availability(self, name: str) -> dict:
        """Check if a company name is available."""
        companies = await self.search_by_name(name, limit=5)

        exact_match = any(c.name.lower() == name.lower() for c in companies)

        return {
            "available": not exact_match,
            "similar_names": [c.name for c in companies],
        }

    async def get_directors(self, registration_number: str) -> list[Director]:
        """Get directors for a company."""
        company = await self.lookup(registration_number)
        return company.directors if company else []

    def is_valid_registration_number(self, reg_number: str) -> bool:
        """Validate registration number format."""
        patterns = [
            r"^\d{4}/\d{6}/\d{2}$",  # 2020/123456/07
            r"^\d{12}$",             # 202012345607
            r"^[A-Z]{2}\d{6}$",      # CK123456 (old CC format)
        ]
        return any(re.match(p, reg_number) for p in patterns)

    def parse_company_type(self, reg_number: str) -> CompanyType | None:
        """Parse company type from registration number."""
        match = re.search(r"/(\d{2})$", reg_number)
        if not match:
            return None

        type_code = match.group(1)
        type_map: dict[str, CompanyType] = {
            "07": "PTY_LTD",
            "06": "LTD",
            "08": "NPC",
            "23": "CC",
            "21": "INC",
        }
        return type_map.get(type_code)

    def format_registration_number(self, reg_number: str) -> str:
        """Format registration number for display."""
        cleaned = re.sub(r"[\s/]", "", reg_number)

        if len(cleaned) == 12:
            return f"{cleaned[:4]}/{cleaned[4:10]}/{cleaned[10:]}"

        return reg_number
