"""Luna SDK - South African Address Utilities for Python."""

from __future__ import annotations

import re
from .types import SAAddress, SAProvince


# Province postal code ranges
POSTAL_CODE_RANGES: dict[SAProvince, list[tuple[int, int]]] = {
    "EC": [(5000, 5999), (6000, 6499)],
    "FS": [(9300, 9999)],
    "GP": [(1, 299), (1000, 2199)],
    "KZN": [(3000, 4499)],
    "LP": [(500, 999)],
    "MP": [(1100, 1399), (2200, 2499)],
    "NC": [(8000, 8999)],
    "NW": [(2500, 2999), (300, 499)],
    "WC": [(6500, 8099)],
}

# SA cities to province mapping
SA_CITIES: dict[str, SAProvince] = {
    "johannesburg": "GP",
    "pretoria": "GP",
    "cape town": "WC",
    "durban": "KZN",
    "port elizabeth": "EC",
    "gqeberha": "EC",
    "bloemfontein": "FS",
    "polokwane": "LP",
    "nelspruit": "MP",
    "mbombela": "MP",
    "kimberley": "NC",
    "rustenburg": "NW",
    "soweto": "GP",
    "sandton": "GP",
    "centurion": "GP",
}

# Province names
PROVINCE_NAMES: dict[SAProvince, str] = {
    "EC": "Eastern Cape",
    "FS": "Free State",
    "GP": "Gauteng",
    "KZN": "KwaZulu-Natal",
    "LP": "Limpopo",
    "MP": "Mpumalanga",
    "NC": "Northern Cape",
    "NW": "North West",
    "WC": "Western Cape",
}


class AddressUtils:
    """South African address utilities."""

    def validate(self, address: SAAddress) -> dict:
        """Validate a South African address."""
        errors: list[str] = []
        warnings: list[str] = []

        # Validate postal code
        if address.postal_code:
            if not re.match(r"^\d{4}$", address.postal_code):
                errors.append("Postal code must be 4 digits")
            else:
                code = int(address.postal_code)
                detected_province = self.get_province_from_postal_code(address.postal_code)
                if address.province and detected_province and address.province != detected_province:
                    warnings.append(
                        f"Postal code {address.postal_code} belongs to {PROVINCE_NAMES.get(detected_province, detected_province)}, "
                        f"but address shows {PROVINCE_NAMES.get(address.province, address.province)}"
                    )

        # Validate province
        if address.province and address.province not in PROVINCE_NAMES:
            errors.append(f"Invalid province code: {address.province}")

        return {
            "valid": len(errors) == 0,
            "errors": errors,
            "warnings": warnings,
        }

    def format(self, address: SAAddress, multiline: bool = True) -> str:
        """Format address for display."""
        parts: list[str] = []

        if address.street:
            parts.append(address.street)
        if address.suburb:
            parts.append(address.suburb)
        if address.city:
            parts.append(address.city)
        if address.province:
            province_name = PROVINCE_NAMES.get(address.province, address.province)
            parts.append(province_name)
        if address.postal_code:
            parts.append(address.postal_code)
        parts.append(address.country)

        separator = "\n" if multiline else ", "
        return separator.join(parts)

    def get_province_from_postal_code(self, postal_code: str) -> SAProvince | None:
        """Detect province from postal code."""
        if not re.match(r"^\d{4}$", postal_code):
            return None

        code = int(postal_code)

        for province, ranges in POSTAL_CODE_RANGES.items():
            for range_start, range_end in ranges:
                if range_start <= code <= range_end:
                    return province

        return None

    def get_province_name(self, code: SAProvince) -> str:
        """Get full province name from code."""
        return PROVINCE_NAMES.get(code, code)

    def lookup_postal_code(self, postal_code: str) -> dict | None:
        """Look up postal code information."""
        if not re.match(r"^\d{4}$", postal_code):
            return None

        province = self.get_province_from_postal_code(postal_code)
        if not province:
            return None

        return {
            "postal_code": postal_code,
            "province": province,
            "province_name": PROVINCE_NAMES.get(province, province),
        }

    def parse(self, address_string: str) -> SAAddress:
        """Parse a free-form address string."""
        result = SAAddress(country="ZA")

        # Try to extract postal code
        postal_match = re.search(r"\b(\d{4})\b", address_string)
        if postal_match and postal_match.group(1):
            result.postal_code = postal_match.group(1)
            result.province = self.get_province_from_postal_code(postal_match.group(1))

        # Try to extract city
        for city, province in SA_CITIES.items():
            if city in address_string.lower():
                result.city = city.title()
                result.province = result.province or province
                break

        return result

    def title_case(self, text: str) -> str:
        """Convert text to title case."""
        return " ".join(word.capitalize() for word in text.split())
